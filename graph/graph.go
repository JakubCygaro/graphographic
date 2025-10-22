package graph

import (
	"container/list"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Node struct {
	Position rl.Vector2
	Contents string
	Edges *list.List
}

type Edge struct {
	Tail *Node
	Head *Node
}

type Graph struct {
	Nodes *list.List
	Edges *list.List
}

func New() Graph {
	return Graph{
		Nodes: list.New(),
		Edges: list.New(),
	}
}
func NewNode() Node {
	return Node {
		Position: rl.Vector2Zero(),
		Contents: "",
		Edges: list.New(),
	}
}
