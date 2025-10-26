package history

import (
	gr "graphographic/graph"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AddNode struct {
	N *gr.Node
}
type RemoveNode struct {
	N *gr.Node
}
type AddEdge struct {
	E *gr.Edge
}
type RemoveEdge struct {
	E *gr.Edge
}
type EditEdgeCost struct {
	E *gr.Edge
	CostPreChange int32
}
type EditNodeContent struct {
	N *gr.Node
	ContentPreChange string
}
type MoveNode struct {
	N *gr.Node
	PosPreChange rl.Vector2
}
