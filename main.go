package main

import (
	"container/list"
	"fmt"
	gr "graphographic/graph"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

const (
	FONT_SIZE      = 24
	FONT_SPACING   = 6
	SCALE_MINIMUM  = 10. / float32(FONT_SIZE)
	LINE_THICKNESS = 4.0
	TARGET_FPS     = 60
)
const (
	MODE_PLACE   = iota
	MODE_APPEND  = iota
	MODE_CONNECT = iota
	MODE_EDIT    = iota
	MODE_MOVE    = iota
	MODE_DELETE  = iota
)

var (
	Width                        = 800
	Height                       = 600
	Graph                        = gr.New()
	Scale             float32    = 1.0
	Offset            rl.Vector2 = rl.Vector2Zero()
	Center            rl.Vector2 = rl.Vector2Scale(rl.Vector2{X: float32(Width), Y: float32(Height)}, 0.5)
	BackgroundColor              = rl.White
	GraphColor                   = rl.Black
	SelectedNodeColor            = rl.SkyBlue
	Mode              int32      = MODE_PLACE
	NodeA             *gr.Node   = nil
	NodeB             *gr.Node   = nil
	Directed          bool       = false
	GridGrain         float32    = 8
	GridSpacing       float32    = float32(Width) / GridGrain
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
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(BackgroundColor)
		update()
		draw()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

// transforms mouse coordinates from screem space to world space
func getMouseWorldPos() rl.Vector2 {
	mousePos := rl.GetMousePosition()
	mousePos = rl.Vector2Subtract(mousePos, Center)
	mousePos = rl.Vector2Subtract(mousePos, Offset)
	mousePos = rl.Vector2Scale(mousePos, 1/Scale)
	return mousePos
}

// transforms coordinates from world space to screen space
func getScreenPos(position rl.Vector2) rl.Vector2 {
	position = rl.Vector2Scale(position, Scale)
	position = rl.Vector2Add(position, Center)
	position = rl.Vector2Add(position, Offset)
	return position
}

// mousePos must be in screen space, node positions will be transformed into screen space
func findNodeUnderMouse(mousePos rl.Vector2) *gr.Node {
	var ret *gr.Node = nil
	for nodeIt := Graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		node := nodeIt.Value.(*gr.Node)
		radius := rl.MeasureTextEx(rl.GetFontDefault(), node.Contents, float32(FONT_SIZE*Scale), FONT_SPACING).X * 0.5 * Scale
		if rl.Vector2Distance(getScreenPos(node.Position), mousePos) <= radius+LINE_THICKNESS {
			return node
		}
	}
	return ret
}
// mousePos must be in screen space, node positions will be transformed into screen space
func findEdgeUnderMouse(mousePos rl.Vector2) *gr.Edge {
	distFromLine := func (p_1, p_2, x rl.Vector2) float32 {
		numerator := ( p_2.Y - p_1.Y ) * x.X - ( p_2.X - p_1.X ) * x.Y + p_2.X * p_1.Y - p_2.Y * p_1.X
		numerator = float32(math.Abs(float64(numerator)))
		denominator := rl.Vector2Distance(p_1, p_2)
		return numerator / denominator
	}

	var ret *gr.Edge = nil
	for edgeIt := Graph.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*gr.Edge)
		dist := distFromLine(edge.StartPos, edge.EndPos, mousePos)
		// fmt.Println("distance: ", dist)
		if dist < 10 {
			return edge
		}
	}
	return ret
}

func update() {
	if rl.IsWindowResized() {
		Width = rl.GetScreenWidth()
		Height = rl.GetScreenHeight()
	}
	mousePos := getMouseWorldPos()
	if Mode == MODE_EDIT && NodeA != nil {
		editModeTyping()
	} else {
		if rl.IsKeyReleased(rl.KeyS) {
			Mode = wrap(Mode+1, int32(MODE_PLACE), int32(MODE_DELETE))
		}
		if rl.IsKeyReleased(rl.KeyE) {
			Mode = MODE_EDIT
		}
		if rl.IsKeyReleased(rl.KeyM) {
			Mode = MODE_MOVE
		}
		if rl.IsKeyReleased(rl.KeyA) {
			Mode = MODE_APPEND
		}
		if rl.IsKeyReleased(rl.KeyP) {
			Mode = MODE_PLACE
		}
		if rl.IsKeyReleased(rl.KeyC) {
			Mode = MODE_CONNECT
		}
		if rl.IsKeyReleased(rl.KeyD) {
			Mode = MODE_DELETE
		}
	}
	if rl.IsKeyReleased(rl.KeyEscape) {
		NodeA = nil
	}

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		switch Mode {
		case MODE_CONNECT:
			if NodeA == nil {
				NodeA = findNodeUnderMouse(rl.GetMousePosition())
			}
		case MODE_APPEND:
			if NodeA == nil {
				NodeA = findNodeUnderMouse(rl.GetMousePosition())
			}
			if NodeA != nil && NodeB == nil {
				NodeB = &gr.Node{
					Position: NodeA.Position,
					Contents: "Node",
					Edges:    list.New(),
				}
			}
		case MODE_MOVE:
			if NodeA == nil {
				NodeA = findNodeUnderMouse(rl.GetMousePosition())
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
			node.Position = mousePos
			node.Contents = "Node"
			Graph.Nodes.PushBack(
				&node,
			)
		case MODE_CONNECT:
			NodeB = findNodeUnderMouse(rl.GetMousePosition())
			if NodeA != nil && NodeB != nil && !NodeA.IsConnectedTo(NodeB) {
				Graph.AddEdge(NodeA, NodeB)
				if !Directed {
					Graph.AddEdge(NodeB, NodeA)
				}
			}
			NodeA = nil
			NodeB = nil
		case MODE_APPEND:
			if NodeA != nil && NodeB != nil {
				Graph.Nodes.PushBack(NodeB)
				Graph.AddEdge(NodeA, NodeB)
				if !Directed {
					Graph.AddEdge(NodeB, NodeA)
				}
			}
			NodeA = nil
			NodeB = nil
		case MODE_EDIT:
			NodeA = findNodeUnderMouse(rl.GetMousePosition())
		case MODE_MOVE:
			NodeA = nil
		case MODE_DELETE:
			if toDelete := findNodeUnderMouse(rl.GetMousePosition()); toDelete != nil {
				Graph.RemoveNode(toDelete)
			}
			if toDelete := findEdgeUnderMouse(rl.GetMousePosition()); toDelete != nil {
				Graph.RemoveEdge(toDelete)
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
		NodeB.Position = mousePos
	}
}

func editModeTyping() {
	selected := NodeA
	if rl.IsKeyPressed(rl.KeyBackspace) {
		end := clamp(len(selected.Contents)-1, 0, len(selected.Contents))
		selected.Contents = selected.Contents[0:end]
	}
	for ch := rl.GetCharPressed(); ch != 0; ch = rl.GetCharPressed() {
		r := rune(ch)
		selected.Contents += string(r)
	}
}

func draw() {
	drawGrid()
	if (Mode == MODE_CONNECT || Mode == MODE_APPEND) && NodeA != nil {
		drawArrow(getScreenPos(NodeA.Position), rl.GetMousePosition(), 15, 10)
		if Mode == MODE_APPEND && NodeB != nil {
			drawNode(NodeB)
		}
	}
	drawGraph()
	var mode string = "Mode: "
	switch Mode {
	case MODE_PLACE:
		mode += "PLACE"
	case MODE_CONNECT:
		mode += "CONNECT"
	case MODE_APPEND:
		mode += "APPEND"
	case MODE_EDIT:
		mode += "EDIT"
	case MODE_MOVE:
		mode += "MOVE"
	case MODE_DELETE:
		mode += "DELETE"
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
	dir = rl.Vector2Rotate(dir, 10 * (math.Pi / 180))
	headPos = rl.Vector2Add(headPos, dir)

	dir = rl.Vector2Subtract(headPos, tailPos)
	dir = rl.Vector2Normalize(dir)
	dir = rl.Vector2Scale(dir, edge.Tail.Radius)
	dir = rl.Vector2Rotate(dir, -10 * (math.Pi / 180))
	tailPos = rl.Vector2Add(tailPos, dir)

	dir = rl.Vector2Subtract(headPos, tailPos)
	halfWay := rl.Vector2Scale(dir, 0.5)

	drawArrow(tailPos, headPos, 15*Scale, 10)
	edge.StartPos = tailPos
	edge.EndPos = headPos

	textPos := rl.Vector2Add(tailPos, halfWay)
	costText := fmt.Sprintf("Cost: %d", edge.Cost)
	size := rl.MeasureTextEx(rl.GetFontDefault(), costText, (FONT_SIZE - 8) * Scale, FONT_SPACING)
	textPos = rl.Vector2Add(textPos, rl.Vector2Scale(rl.Vector2Rotate(rl.Vector2Normalize(halfWay), 90 * math.Pi / 180), -20))

	rot := rl.Vector2Angle(rl.Vector2Normalize(halfWay), rl.Vector2 { X: 1, Y: 0 })
	rl.DrawTextPro(
		rl.GetFontDefault(),
		costText,
		textPos,
		rl.Vector2Scale(size, 0.5),
		(-(rot / ( math.Pi / 180 )) ),
		(FONT_SIZE - 8) * Scale,
		FONT_SPACING,
		rl.Red,
	)

}
func drawNode(node *gr.Node) {
	radius := rl.MeasureTextEx(rl.GetFontDefault(), node.Contents, float32(FONT_SIZE*Scale), FONT_SPACING).X * 0.5 * Scale
	radius *= 1.2
	position := getScreenPos(node.Position)
	textPosition := position
	var color rl.Color
	amISelected := (Mode == MODE_EDIT || Mode == MODE_MOVE) && node == NodeA
	if amISelected {
		color = SelectedNodeColor
	} else {
		color = GraphColor
	}
	node.Radius = radius + LINE_THICKNESS
	rl.DrawCircle(int32(position.X), int32(position.Y), radius+LINE_THICKNESS, color)

	rl.DrawCircle(int32(position.X), int32(position.Y), radius, BackgroundColor)

	text := node.Contents
	size := rl.MeasureTextEx(rl.GetFontDefault(), node.Contents, FONT_SIZE*Scale, FONT_SPACING-2)
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
}

func drawArrow(a, b rl.Vector2, h, w float32) {
	dir := rl.Vector2Subtract(a, b)
	dir = rl.Vector2Normalize(dir)
	height := rl.Vector2Scale(dir, h)
	b_ := rl.Vector2Add(b, height)

	rl.DrawLineEx((a), (b_), LINE_THICKNESS, GraphColor)

	width := rl.Vector2Rotate(dir, -90*(math.Pi/180))
	width = rl.Vector2Scale(width, w)

	x := rl.Vector2Add(b, width)
	x = rl.Vector2Add(x, height)

	y := b

	z := rl.Vector2Subtract(b, width)
	z = rl.Vector2Add(z, height)

	rl.DrawTriangle((x), (y), (z), GraphColor)

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
