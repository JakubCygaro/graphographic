package main

import (
	"container/list"
	"fmt"
	algo "graphographic/algorithm"
	gr "graphographic/graph"
	hist "graphographic/history"
	"math"
	"strconv"
	"unicode"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

const (
	FONT_SIZE               = 24
	FONT_SPACING            = 6
	SCALE_MINIMUM           = 10. / float32(FONT_SIZE)
	LINE_THICKNESS          = 4.0
	TARGET_FPS              = 60
	ACTION_HISTORY_MAX_SIZE = 40
	MIN_RADIUS              = 10
)
const (
	MODE_PLACE     = iota
	MODE_APPEND    = iota
	MODE_CONNECT   = iota
	MODE_EDIT      = iota
	MODE_MOVE      = iota
	MODE_ALGORITHM = iota
	MODE_DELETE    = iota
)

var (
	Width                          = 800
	Height                         = 600
	Graph                          = gr.New()
	Scale               float32    = 1.0
	Offset              rl.Vector2 = rl.Vector2Zero()
	Center              rl.Vector2 = rl.Vector2Scale(rl.Vector2{X: float32(Width), Y: float32(Height)}, 0.5)
	BackgroundColor                = rl.White
	GraphColor                     = rl.Black
	SelectedNodeColor              = rl.SkyBlue
	Mode                int32      = MODE_PLACE
	NodeA               *gr.Node   = nil
	NodeB               *gr.Node   = nil
	EdgeA               *gr.Edge   = nil
	SelectedEdgeScratch string     = ""
	Directed            bool       = false
	GridGrain           float32    = 12
	GridSpacing         float32    = float32(Width) / GridGrain
	ActionHistory       []any      = make([]any, 0)
	// mouse position in screen space
	MousePos             rl.Vector2
	Algorithms           []algo.Algorithm = make([]algo.Algorithm, 0)
	CurrentAlgorithm     int              = 0
	CurrentAlgorithmName string           = "unnamed"
	IsAlgorithmRunning   bool             = false
	AlgorithmSpeed       int              = 30
	AlgorithmErrorMsg    string           = ""

	UpdateCounter uint64 = 0
)

func clamp[T constraints.Ordered](value, min, max T) T {
	if value > max {
		return max
	} else if value < min {
		return min
	} else {
		return value
	}
}

func wrap[T constraints.Ordered](value, min, max T) T {
	if value > max {
		return min
	} else if value < min {
		return max
	} else {
		return value
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(Width), int32(Height), "Graphographic")
	rl.SetTargetFPS(TARGET_FPS)
	initGraph()
	loadAlgorithms()
	spreadNodes()
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(BackgroundColor)
		update()
		draw()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

func initGraph() {
	node := gr.NewNode()
	node.Content = "Node A"
	a := Graph.AddNode(node)

	node = gr.NewNode()
	node.Content = "Node B"
	b := Graph.AddNode(node)

	Graph.AddEdge(a, b).Cost = 10

	node = gr.NewNode()
	node.Content = "Node C"
	c := Graph.AddNode(node)

	Graph.AddEdge(a, c).Cost = 3

	node = gr.NewNode()
	node.Content = "Node D"
	d := Graph.AddNode(node)

	Graph.AddEdge(c, d).Cost = 12

	node = gr.NewNode()
	node.Content = "Node F"
	f := Graph.AddNode(node)

	Graph.AddEdge(c, f).Cost = 8
	Graph.AddEdge(d, b).Cost = 7
}

func loadAlgorithms() {
	dfs := &algo.DFS{}
	Algorithms = append(Algorithms, dfs)
	bfs := &algo.BFS{}
	Algorithms = append(Algorithms, bfs)
	dijkstra := &algo.Dijkstra{}
	Algorithms = append(Algorithms, dijkstra)

	CurrentAlgorithmName = Algorithms[CurrentAlgorithm].GetName()
}

func resetAlgoDataState() {
	for nodeIt := Graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		n := nodeIt.Value.(*gr.Node)
		n.Data = gr.AlgoData{}
	}
	for edgeIt := Graph.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		e := edgeIt.Value.(*gr.Edge)
		e.Data = gr.AlgoData{}
	}
}

func spreadNodes() {
	type pair struct {
		n        *gr.Node
		displace rl.Vector2
	}
	nodes := make([]pair, 0)
	for nodeIt := Graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		node := nodeIt.Value.(*gr.Node)
		node.Radius = calculateNodeRadius(node)
		nodes = append(nodes, pair{
			n:        node,
			displace: rl.Vector2Zero(),
		})
	}
	i := 10
	for hadOverlap := true; hadOverlap; {
		i--
		hadOverlap = false
		for i := 0; i < len(nodes); i++ {
			for j := i + 1; j < len(nodes); j++ {
				a, b := &nodes[i], &nodes[j]
				if rl.Vector2Distance(a.n.Position, b.n.Position) <= a.n.Radius+b.n.Radius {
					hadOverlap = true
					dir := rl.Vector2Subtract(
						a.n.Position,
						b.n.Position,
					)
					if dir == rl.Vector2Zero() {
						dir = rl.Vector2{X: 1, Y: 0}
						dir = rl.Vector2Rotate(dir, float32(i)/float32(len(nodes))*180./math.Pi)
					}
					dir = rl.Vector2Scale(dir, a.n.Radius*2)
					a.displace = rl.Vector2Add(
						a.displace,
						dir,
					)
					b.displace = rl.Vector2Add(
						a.displace,
						rl.Vector2Scale(dir, -2),
					)
				}
			}
			a := &nodes[i]
			a.n.Position = rl.Vector2Add(
				a.n.Position,
				a.displace,
			)
			a.displace = rl.Vector2Zero()
		}
	}
}

// transforms mouse coordinates from screem space to world space
func getMouseWorldPos() rl.Vector2 {
	mousePos := MousePos
	mousePos = rl.Vector2Subtract(mousePos, Center)
	mousePos = rl.Vector2Subtract(mousePos, Offset)
	mousePos = rl.Vector2Scale(mousePos, 1/Scale)
	return mousePos
}

// transforms coordinates from world space to screen space
func getScreenPos(position rl.Vector2) rl.Vector2 {
	position = rl.Vector2Add(position, Center)
	position = rl.Vector2Add(position, Offset)
	position = rl.Vector2Scale(position, Scale)
	return position
}

// mousePos must be in screen space, node positions will be transformed into screen space
func findNodeUnderMouse() *gr.Node {
	var ret *gr.Node = nil
	for nodeIt := Graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		node := nodeIt.Value.(*gr.Node)
		if isNodeUnderMouse(node) {
			return node
		}
	}
	return ret
}
func isNodeUnderMouse(node *gr.Node) bool {
	radius := node.Radius
	return rl.Vector2Distance(getScreenPos(node.Position), MousePos) <= radius
}

// mousePos must be in screen space, node positions will be transformed into screen space
func findEdgeUnderMouse() *gr.Edge {
	var ret *gr.Edge = nil
	for edgeIt := Graph.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*gr.Edge)
		if isEdgeUnderMouse(edge) {
			return edge
		}
	}
	return ret
}
func isEdgeUnderMouse(edge *gr.Edge) bool {
	distFromLine := func(p_1, p_2, x rl.Vector2) float32 {
		numerator := (p_2.Y-p_1.Y)*x.X - (p_2.X-p_1.X)*x.Y + p_2.X*p_1.Y - p_2.Y*p_1.X
		numerator = float32(math.Abs(float64(numerator)))
		denominator := rl.Vector2Distance(p_1, p_2)
		res := numerator / denominator
		return res
	}
	dir := rl.Vector2Subtract(edge.EndPos, edge.StartPos)
	halfWay := rl.Vector2Scale(dir, 0.5)
	center := rl.Vector2Add(edge.StartPos, halfWay)
	distFromCenter := rl.Vector2Distance(center, MousePos)
	return distFromLine(edge.StartPos, edge.EndPos, MousePos) < 10 && distFromCenter < rl.Vector2Length(halfWay)+10
}

func revertLatestAction() {
	if len(ActionHistory) == 0 {
		return
	}
	latest := ActionHistory[len(ActionHistory)-1]
	ActionHistory = ActionHistory[0 : len(ActionHistory)-1]
	if addN, ok := latest.(*hist.AddNode); ok {
		Graph.RemoveNode(addN.N)
	} else if addE, ok := latest.(*hist.AddEdge); ok {
		Graph.RemoveEdge(addE.E)
	} else if remN, ok := latest.(*hist.RemoveNode); ok {
		Graph.Nodes.PushBack(remN.N)
	} else if remE, ok := latest.(*hist.RemoveEdge); ok {
		Graph.Edges.PushBack(remE.E)
		remE.E.Head.Edges.PushBack(remE.E)
		remE.E.Tail.Edges.PushBack(remE.E)
	} else if chE, ok := latest.(*hist.EditEdgeCost); ok {
		chE.E.Cost = chE.CostPreChange
	} else if chN, ok := latest.(*hist.EditNodeContent); ok {
		chN.N.Content = chN.ContentPreChange
	} else if chN, ok := latest.(*hist.MoveNode); ok {
		chN.N.Position = chN.PosPreChange
	} else if _, ok := latest.(*hist.NodeSelected); ok {
		Algorithms[CurrentAlgorithm].UndoSelect()
	}

}

func update() {
	MousePos = rl.GetMousePosition()
	UpdateCounter++
	mousePosWorld := getMouseWorldPos()

	if IsAlgorithmRunning && UpdateCounter%uint64(AlgorithmSpeed) == 0 {
		IsAlgorithmRunning = Algorithms[CurrentAlgorithm].Update()
	}

	if rl.IsWindowResized() {
		Width = rl.GetScreenWidth()
		Height = rl.GetScreenHeight()
	}
	if Mode == MODE_EDIT && NodeA != nil || EdgeA != nil {
		editModeTyping()
	} else {
		if rl.IsKeyReleased(rl.KeyE) {
			Mode = MODE_EDIT
		}
		if rl.IsKeyReleased(rl.KeyM) {
			Mode = MODE_MOVE
		}
		if rl.IsKeyReleased(rl.KeyA) {
			Mode = MODE_APPEND
		}
		if rl.IsKeyReleased(rl.KeyT) {
			Mode = MODE_ALGORITHM
			resetAlgoDataState()
			Algorithms[CurrentAlgorithm].Init()
		}
		if rl.IsKeyReleased(rl.KeyP) {
			Mode = MODE_PLACE
		}
		if rl.IsKeyReleased(rl.KeyC) {
			Mode = MODE_CONNECT
		}
		if rl.IsKeyReleased(rl.KeyD) && rl.IsKeyDown(rl.KeyLeftShift) {
			Directed = !Directed
		} else if rl.IsKeyReleased(rl.KeyD) {
			Mode = MODE_DELETE
		}
		if rl.IsKeyReleased(rl.KeyBackspace) {
			revertLatestAction()
		}
		if rl.IsKeyReleased(rl.KeyS) {
			CurrentAlgorithm = wrap(CurrentAlgorithm+1, 0, len(Algorithms)-1)
			resetAlgoDataState()
			Algorithms[CurrentAlgorithm].Init()
			CurrentAlgorithmName = Algorithms[CurrentAlgorithm].GetName()
		}
		if rl.IsKeyReleased(rl.KeyR) && Mode == MODE_ALGORITHM {
			resetAlgoDataState()
			if err := Algorithms[CurrentAlgorithm].Start(&Graph); err != nil {
				rl.TraceLog(rl.LogWarning, "%s", err.Error())
				IsAlgorithmRunning = false
				AlgorithmErrorMsg = err.Error()
			} else {
				IsAlgorithmRunning = true
			}
		}
	}
	if rl.IsKeyReleased(rl.KeyEscape) {
		NodeA = nil
	}

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		switch Mode {
		case MODE_CONNECT:
			if NodeA == nil {
				NodeA = findNodeUnderMouse()
			}
		case MODE_APPEND:
			if NodeA == nil {
				NodeA = findNodeUnderMouse()
			}
			if NodeA != nil && NodeB == nil {
				NodeB = &gr.Node{
					Position: NodeA.Position,
					Content:  "Node",
					Edges:    list.New(),
				}
			}
		case MODE_MOVE:
			if NodeA == nil {
				NodeA = findNodeUnderMouse()
				if NodeA != nil {
					ActionHistory = append(ActionHistory, &hist.MoveNode{N: NodeA, PosPreChange: NodeA.Position})
				}
			}
			if NodeA != nil {
				NodeA.Position = getMouseWorldPos()
			}
		}
	}
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {

		switch Mode {
		case MODE_PLACE:
			node := gr.NewNode()
			node.Position = mousePosWorld
			node.Content = "Node"
			Graph.Nodes.PushBack(
				&node,
			)
			ActionHistory = append(ActionHistory, &hist.AddNode{N: &node})
		case MODE_CONNECT:
			NodeB = findNodeUnderMouse()
			if NodeA != nil && NodeB != nil && !NodeA.IsConnectedTo(NodeB) && NodeA != NodeB {
				edge := Graph.AddEdge(NodeA, NodeB)
				ActionHistory = append(ActionHistory, &hist.AddEdge{E: edge})
				if !Directed {
					edge := Graph.AddEdge(NodeB, NodeA)
					ActionHistory = append(ActionHistory, &hist.AddEdge{E: edge})
				}
			}
			NodeA = nil
			NodeB = nil
		case MODE_APPEND:
			if NodeA != nil && NodeB != nil {
				Graph.Nodes.PushBack(NodeB)
				ActionHistory = append(ActionHistory, &hist.AddNode{N: NodeB})
				edge := Graph.AddEdge(NodeA, NodeB)
				ActionHistory = append(ActionHistory, &hist.AddEdge{E: edge})
				if !Directed {
					edge := Graph.AddEdge(NodeB, NodeA)
					ActionHistory = append(ActionHistory, &hist.AddEdge{E: edge})
				}
			}
			NodeA = nil
			NodeB = nil
		case MODE_EDIT:
			EdgeA = nil
			NodeA = nil
			NodeA = findNodeUnderMouse()
			if NodeA == nil {
				if EdgeA = findEdgeUnderMouse(); EdgeA != nil {
					ActionHistory = append(ActionHistory, &hist.EditEdgeCost{E: EdgeA, CostPreChange: EdgeA.Cost})
					SelectedEdgeScratch = fmt.Sprintf("%d", EdgeA.Cost)
				}
			} else {
				ActionHistory = append(ActionHistory, &hist.EditNodeContent{N: NodeA, ContentPreChange: NodeA.Content})
			}
		case MODE_MOVE:
			NodeA = nil
		case MODE_ALGORITHM:
			slc := findNodeUnderMouse()
			Algorithms[CurrentAlgorithm].NodeSelected(slc)
			if slc != nil {
				ActionHistory = append(ActionHistory, &hist.NodeSelected{N: slc})
			}
		case MODE_DELETE:
			if toDelete := findNodeUnderMouse(); toDelete != nil {
				Graph.RemoveNode(toDelete)
				ActionHistory = append(ActionHistory, &hist.RemoveNode{N: toDelete})
			}
			if toDelete := findEdgeUnderMouse(); toDelete != nil {
				Graph.RemoveEdge(toDelete)
				ActionHistory = append(ActionHistory, &hist.RemoveEdge{E: toDelete})
			}
		}
	}

	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		delta := rl.GetMouseDelta()
		Offset = rl.Vector2Add(delta, Offset)
	}

	scrollY := rl.GetMouseWheelMoveV().Y * 0.01
	if scrollY != 0.0 {
		Scale += scrollY
		Scale = rl.Clamp(Scale, SCALE_MINIMUM, 100)
	}
	if Mode == MODE_APPEND && NodeB != nil {
		NodeB.Position = mousePosWorld
	}
	if len(ActionHistory) > ACTION_HISTORY_MAX_SIZE {
		_, ActionHistory = ActionHistory[0], ActionHistory[1:]
	}
}

func editModeTyping() {
	if rl.IsKeyPressed(rl.KeyBackspace) {
		if selected := NodeA; selected != nil {
			end := clamp(len(selected.Content)-1, 0, len(selected.Content))
			selected.Content = selected.Content[0:end]
		} else if selected := EdgeA; selected != nil {
			if len(SelectedEdgeScratch) == 1 {
				SelectedEdgeScratch = "0"
			} else {
				end := clamp(len(SelectedEdgeScratch)-1, 0, len(SelectedEdgeScratch))
				SelectedEdgeScratch = SelectedEdgeScratch[0:end]
			}
			if d, err := strconv.Atoi(SelectedEdgeScratch); err == nil {
				selected.Cost = int32(d)
			} else {
				SelectedEdgeScratch = fmt.Sprintf("%d", selected.Cost)
			}
		}
	}
	for ch := rl.GetCharPressed(); ch != 0; ch = rl.GetCharPressed() {
		r := rune(ch)
		if selected := NodeA; selected != nil {
			selected.Content += string(r)
		} else if selected := EdgeA; selected != nil && (unicode.IsDigit(r) || r == '-') {
			if SelectedEdgeScratch == "0" {
				SelectedEdgeScratch = string(r)
			} else if r == '-' {
				SelectedEdgeScratch = strconv.Itoa(int(-selected.Cost))
			} else {
				SelectedEdgeScratch += string(r)
			}

			if d, err := strconv.Atoi(SelectedEdgeScratch); err == nil {
				selected.Cost = int32(d)
			} else {
				SelectedEdgeScratch = fmt.Sprintf("%d", selected.Cost)
			}
		}
	}
}

func draw() {
	drawGrid()
	if (Mode == MODE_CONNECT || Mode == MODE_APPEND) && NodeA != nil {
		drawArrow(getScreenPos(NodeA.Position), MousePos, 15, 10, GraphColor)
		if Mode == MODE_APPEND && NodeB != nil {
			drawNode(NodeB)
		}
	}
	drawGraph()
	var mode string = "Mode: "
	var directed string
	if Directed {
		directed = "(DIRECTED)"
	} else {
		directed = "(UNDIRECTED)"
	}
	switch Mode {
	case MODE_PLACE:
		mode += "PLACE"
	case MODE_CONNECT:
		mode += "CONNECT " + directed
	case MODE_APPEND:
		mode += "APPEND"
	case MODE_EDIT:
		mode += "EDIT"
	case MODE_MOVE:
		mode += "MOVE"
	case MODE_DELETE:
		mode += "DELETE"
	case MODE_ALGORITHM:
		mode += "ALGORITHM"
	}
	size := rl.MeasureTextEx(rl.GetFontDefault(), mode, FONT_SIZE, FONT_SPACING)
	rl.DrawTextEx(
		rl.GetFontDefault(),
		mode,
		rl.Vector2{X: 0, Y: float32(Height) - size.Y},
		FONT_SIZE,
		FONT_SPACING,
		rl.Red,
	)
	scaleText := fmt.Sprintf("%.3fx", Scale)
	size = rl.MeasureTextEx(rl.GetFontDefault(), scaleText, FONT_SIZE-4, FONT_SPACING)
	rl.DrawTextEx(
		rl.GetFontDefault(),
		scaleText,
		rl.Vector2{X: 0, Y: 0},
		FONT_SIZE,
		FONT_SPACING,
		rl.Red,
	)
	algoName := "Algorithm: " + CurrentAlgorithmName
	size = rl.MeasureTextEx(rl.GetFontDefault(), algoName, FONT_SIZE-4, FONT_SPACING)
	rl.DrawTextEx(
		rl.GetFontDefault(),
		algoName,
		rl.Vector2{X: float32(Width) - size.X, Y: 0},
		FONT_SIZE-4,
		FONT_SPACING,
		rl.Red,
	)
	if AlgorithmErrorMsg != "" {
		if !IsAlgorithmRunning {
			algoErr := "Error: " + AlgorithmErrorMsg
			size = rl.MeasureTextEx(rl.GetFontDefault(), algoErr, FONT_SIZE-6, FONT_SPACING)
			rl.DrawTextEx(
				rl.GetFontDefault(),
				algoErr,
				rl.Vector2{X: float32(Width) - size.X, Y: float32(Height) - size.Y},
				FONT_SIZE-6,
				FONT_SPACING,
				rl.Green,
			)

		} else {
			AlgorithmErrorMsg = ""
		}
	}
}

func drawGraph() {
	// draw edges
	for edgeIt := Graph.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*gr.Edge)
		drawEdge(edge)
	}
	// draw nodes
	for nodeIt := Graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		node := nodeIt.Value.(*gr.Node)
		drawNode(node)
	}
}
func drawEdge(edge *gr.Edge) {
	tailPos := edge.Tail.Position
	headPos := edge.Head.Position

	tailPos = getScreenPos(tailPos)
	headPos = getScreenPos(headPos)

	dir := rl.Vector2Subtract(tailPos, headPos)
	dir = rl.Vector2Normalize(dir)
	dir = rl.Vector2Scale(dir, edge.Head.Radius)
	dir = rl.Vector2Rotate(dir, 10*(math.Pi/180))
	headPos = rl.Vector2Add(headPos, dir)

	dir = rl.Vector2Subtract(headPos, tailPos)
	dir = rl.Vector2Normalize(dir)
	dir = rl.Vector2Scale(dir, edge.Tail.Radius)
	dir = rl.Vector2Rotate(dir, -10*(math.Pi/180))
	tailPos = rl.Vector2Add(tailPos, dir)

	dir = rl.Vector2Subtract(headPos, tailPos)
	halfWay := rl.Vector2Scale(dir, 0.5)

	var color rl.Color

	if edge.Data.Highlighted && Mode == MODE_ALGORITHM {
		color = SelectedNodeColor
	} else if edge.Data.Explored && Mode == MODE_ALGORITHM {
		color = rl.Green
	} else if edge == EdgeA && Mode == MODE_EDIT {
		color = SelectedNodeColor
	} else if Mode == MODE_DELETE && isEdgeUnderMouse(edge) {
		color = rl.Red
	} else {
		color = GraphColor
	}

	drawArrow(tailPos, headPos, 15*Scale, 10, color)
	edge.StartPos = tailPos
	edge.EndPos = headPos

	textPos := rl.Vector2Add(tailPos, halfWay)
	costText := fmt.Sprintf("Cost: %d", edge.Cost)
	size := rl.MeasureTextEx(rl.GetFontDefault(), costText, (FONT_SIZE-8)*Scale, FONT_SPACING)
	textPos = rl.Vector2Add(textPos, rl.Vector2Scale(rl.Vector2Rotate(rl.Vector2Normalize(halfWay), 90*math.Pi/180), -20))

	rot := rl.Vector2Angle(rl.Vector2Normalize(halfWay), rl.Vector2{X: 1, Y: 0})
	rl.DrawTextPro(
		rl.GetFontDefault(),
		costText,
		textPos,
		rl.Vector2Scale(size, 0.5),
		(-(rot / (math.Pi / 180))),
		(FONT_SIZE-8)*Scale,
		FONT_SPACING,
		rl.Red,
	)

}

func drawNode(node *gr.Node) {
	radius := calculateNodeRadius(node)
	radius *= 1.2
	position := getScreenPos(node.Position)
	textPosition := position
	var color rl.Color
	amISelected := (Mode == MODE_EDIT || Mode == MODE_MOVE) && node == NodeA
	if node.Data.Highlighted && Mode == MODE_ALGORITHM {
		color = SelectedNodeColor
	} else if node.Data.Explored && Mode == MODE_ALGORITHM {
		color = rl.Green
	} else if amISelected {
		color = SelectedNodeColor
	} else if Mode == MODE_DELETE && isNodeUnderMouse(node) {
		color = rl.Red
	} else {
		color = GraphColor
	}
	node.Radius = radius + LINE_THICKNESS
	rl.DrawCircle(int32(position.X), int32(position.Y), radius+LINE_THICKNESS, color)

	rl.DrawCircle(int32(position.X), int32(position.Y), radius, BackgroundColor)

	text := node.Content
	size := rl.MeasureTextEx(rl.GetFontDefault(), node.Content, FONT_SIZE*Scale, FONT_SPACING-2)
	if size.X > 2*radius && !amISelected {
		text = "..."
		size = rl.MeasureTextEx(rl.GetFontDefault(), text, FONT_SIZE*Scale, FONT_SPACING-2)
	}
	textPosition = rl.Vector2Subtract(textPosition, rl.Vector2Scale(size, 0.5))

	rl.DrawTextEx(
		rl.GetFontDefault(),
		text,
		textPosition,
		FONT_SIZE*Scale,
		FONT_SPACING-2,
		rl.Red,
	)

	if node.Data.Tag != "" && Mode == MODE_ALGORITHM {
		text := node.Data.Tag
		size := rl.MeasureTextEx(rl.GetFontDefault(), text, FONT_SIZE*Scale, FONT_SPACING-2)
		textPosition = rl.Vector2Add(position, rl.Vector2{X: 0, Y: -radius})
		textPosition = rl.Vector2Subtract(textPosition, rl.Vector2Scale(size, 0.5))

		rl.DrawRectangleV(
			textPosition,
			size,
			BackgroundColor,
		)
		rl.DrawTextEx(
			rl.GetFontDefault(),
			text,
			textPosition,
			FONT_SIZE*Scale,
			FONT_SPACING-2,
			rl.Red,
		)
	}
}
func calculateNodeRadius(node *gr.Node) float32 {
	radius := rl.MeasureTextEx(rl.GetFontDefault(), node.Content, float32(FONT_SIZE*Scale), FONT_SPACING).X * 0.5 * Scale
	radius = float32(math.Max(float64(radius), MIN_RADIUS))
	return radius
}
func drawArrow(a, b rl.Vector2, h, w float32, color rl.Color) {
	dir := rl.Vector2Subtract(a, b)
	dir = rl.Vector2Normalize(dir)
	height := rl.Vector2Scale(dir, h)
	b_ := rl.Vector2Add(b, height)

	rl.DrawLineEx((a), (b_), LINE_THICKNESS, color)

	width := rl.Vector2Rotate(dir, -90*(math.Pi/180))
	width = rl.Vector2Scale(width, w)

	x := rl.Vector2Add(b, width)
	x = rl.Vector2Add(x, height)

	y := b

	z := rl.Vector2Subtract(b, width)
	z = rl.Vector2Add(z, height)

	rl.DrawTriangle((x), (y), (z), color)

}

func drawGrid() {
	spacing := GridSpacing * Scale

	gridCenter := rl.Vector2{
		X: float32(Width / 2.),
		Y: float32(Height / 2.),
	}
	// Y axis
	for i := -GridGrain; i <= GridGrain; i++ {
		rl.DrawLineV(
			rl.Vector2{
				X: 0,
				Y: gridCenter.Y + i*spacing,
			},
			rl.Vector2{
				X: float32(Width),
				Y: gridCenter.Y + i*spacing,
			},
			rl.Gray,
		)
	}
	// X axis
	for i := -GridGrain; i <= GridGrain; i++ {
		rl.DrawLineV(
			rl.Vector2{
				X: gridCenter.X + i*spacing,
				Y: 0,
			},
			rl.Vector2{
				X: gridCenter.X + i*spacing,
				Y: float32(Height),
			},
			rl.Gray,
		)
	}
}
