package main

import (
	"container/list"
	gr "graphographic/graph"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

const (
	WIDTH          = 800
	HEIGHT         = 600
	FONT_SIZE      = 24
	FONT_SPACING   = 10
	SCALE_MINIMUM  = 10. / float32(FONT_SIZE)
	LINE_THICKNESS = 4.0
)
const (
	MODE_PLACE   = iota
	MODE_APPEND  = iota
	MODE_CONNECT = iota
)

var (
	Graph                      = gr.New()
	Scale           float32    = 1.0
	Offset          rl.Vector2 = rl.Vector2Zero()
	Center          rl.Vector2 = rl.Vector2Scale(rl.Vector2{X: WIDTH, Y: HEIGHT}, 0.5)
	BackgroundColor            = rl.White
	GraphColor                 = rl.Black
	Mode            int32      = MODE_PLACE
	NodeA           *gr.Node   = nil
	NodeB           *gr.Node   = nil
	Directed        bool       = false
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
	rl.InitWindow(WIDTH, HEIGHT, "Graphographic")

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
		if rl.Vector2Distance(getScreenPos(node.Position), mousePos) <= radius + LINE_THICKNESS {
			return node
		}
	}
	return ret
}

func update() {
	mousePos := getMouseWorldPos()
	if rl.IsKeyReleased(rl.KeyC) {
		Mode = wrap(Mode+1, int32(MODE_PLACE), int32(MODE_CONNECT))
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
					Contents: "New Node",
					Edges: list.New(),
				}
			}
		}
	}
	if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {

		switch Mode {
		case MODE_PLACE:
			node := gr.NewNode()
			node.Position = mousePos
			node.Contents = "New node"
			Graph.Nodes.PushBack(
				&node,
			)
		case MODE_CONNECT:
			NodeB = findNodeUnderMouse(rl.GetMousePosition())
			if NodeA != nil && NodeB != nil && !NodeA.IsConnectedTo(NodeB){
				Graph.AddEdge(NodeA, NodeB)
				if !Directed {
					Graph.AddEdge(NodeB, NodeA)
				}
			}
		case MODE_APPEND:
			if NodeA != nil && NodeB != nil {
				Graph.Nodes.PushBack(NodeB)
				Graph.AddEdge(NodeA, NodeB)
				if !Directed {
					Graph.AddEdge(NodeB, NodeA)
				}
			}
		}
		NodeA = nil
		NodeB = nil
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
		NodeB.Position = mousePos;
	}
}

func draw() {
	if NodeA != nil {
		positionA := getScreenPos(NodeA.Position)
		mousePos := rl.GetMousePosition()
		rl.DrawLineEx(
			mousePos,
			positionA,
			LINE_THICKNESS,
			GraphColor,
		)
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
	}
	size := rl.MeasureTextEx(rl.GetFontDefault(), mode, FONT_SIZE, FONT_SPACING)
	rl.DrawTextEx(
		rl.GetFontDefault(),
		mode,
		rl.Vector2{X: 0, Y: HEIGHT - size.Y},
		FONT_SIZE,
		FONT_SPACING,
		rl.Red,
	)
}

func drawGraph() {
	// draw edges
	for edgeIt := Graph.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*gr.Edge)
		tailPos := getScreenPos(edge.Tail.Position)
		headPos := getScreenPos(edge.Head.Position)
		rl.DrawLineEx(
			tailPos,
			headPos,
			LINE_THICKNESS,
			GraphColor,
		)
	}
	// draw nodes
	for nodeIt := Graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		node := nodeIt.Value.(*gr.Node)
		drawNode(node)
	}
}

func drawNode(node *gr.Node){
	radius := rl.MeasureTextEx(rl.GetFontDefault(), node.Contents, float32(FONT_SIZE*Scale), FONT_SPACING).X * 0.5 * Scale
	position := getScreenPos(node.Position)
	textPosition := position
	rl.DrawCircle(int32(position.X), int32(position.Y), radius+LINE_THICKNESS, GraphColor)
	rl.DrawCircle(int32(position.X), int32(position.Y), radius, BackgroundColor)

	size := rl.MeasureTextEx(rl.GetFontDefault(), node.Contents, FONT_SIZE*Scale, FONT_SPACING)
	textPosition = rl.Vector2Subtract(textPosition, rl.Vector2Scale(size, 0.5))

	rl.DrawTextEx(
		rl.GetFontDefault(),
		node.Contents,
		textPosition,
		FONT_SIZE*Scale,
		FONT_SPACING,
		rl.Red,
	)
}
