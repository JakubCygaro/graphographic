package main

import (
	gr "graphographic/graph"

	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/constraints"
)

const (
	WIDTH         = 800
	HEIGHT        = 600
	FONT_SIZE     = 24
	FONT_SPACING  = 10
	SCALE_MINIMUM = 10. / float32(FONT_SIZE)
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
	Mode            int32      = MODE_APPEND
	NodeA           *gr.Node   = nil
	NodeB           *gr.Node   = nil
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


func findNodeUnderMouse(mousePos rl.Vector2) *gr.Node {
	var ret *gr.Node = nil
	for nodeIt := Graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		node := nodeIt.Value.(*gr.Node)
		radius := rl.MeasureTextEx(rl.GetFontDefault(), node.Contents, float32(FONT_SIZE*Scale), FONT_SPACING).X * 0.5 * Scale
		if rl.Vector2Distance(node.Position, mousePos) <= radius {
			return node
		}
	}
	return ret
}

func getTransformedMousePos() rl.Vector2 {
	mousePos := rl.GetMousePosition()
	mousePos = rl.Vector2Subtract(mousePos, Center)
	mousePos = rl.Vector2Subtract(mousePos, Offset)
	mousePos = rl.Vector2Scale(mousePos, 1/Scale)
	return mousePos
}
func getTransformedPos(position rl.Vector2) rl.Vector2 {
	position = rl.Vector2Scale(position, Scale)
	position = rl.Vector2Add(position, Center)
	position = rl.Vector2Add(position, Offset)
	return position
}

func update() {
	mousePos := getTransformedMousePos()
	if rl.IsKeyReleased(rl.KeyC) {
		Mode = wrap(Mode+1, int32(MODE_PLACE), int32(MODE_CONNECT))
	}

	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		switch Mode {
		case MODE_CONNECT:
			if NodeA == nil {
				NodeA = findNodeUnderMouse(mousePos)
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
			NodeB = findNodeUnderMouse(mousePos)
			if NodeA != nil && NodeB != nil {
				Graph.Edges.PushBack(&gr.Edge {
					Tail: NodeA,
					Head: NodeB,
				})
				Graph.Edges.PushBack(&gr.Edge {
					Head: NodeA,
					Tail: NodeB,
				})
			}
			NodeA = nil
			NodeB = nil
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
}

func draw() {
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
	if NodeA != nil {
		mousePos := rl.GetMousePosition()
		position := getTransformedPos(NodeA.Position)
		rl.DrawLine(
			int32(position.X),
			int32(position.Y),
			int32(mousePos.X),
			int32(mousePos.Y),
			rl.Black,
		)
	}
}


func drawGraph() {
	for edgeIt := Graph.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*gr.Edge)
		tailPos := getTransformedPos(edge.Tail.Position)
		headPos := getTransformedPos(edge.Head.Position)
		rl.DrawLine(
			int32(tailPos.X),
			int32(tailPos.Y),
			int32(headPos.X),
			int32(headPos.Y),
			rl.Black,
		)
	}
	for nodeIt := Graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		node := nodeIt.Value.(*gr.Node)
		radius := rl.MeasureTextEx(rl.GetFontDefault(), node.Contents, float32(FONT_SIZE*Scale), FONT_SPACING).X * 0.5 * Scale
		position := getTransformedPos(node.Position)
		textPosition := position
		rl.DrawCircle(int32(position.X), int32(position.Y), radius+3.0, GraphColor)
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
}
