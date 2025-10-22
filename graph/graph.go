package graph

import (
	"container/list"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Node struct {
	Position rl.Vector2
	Contents string
	Edges *list.List
	// saved from draw pass
	Radius float32
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

// Connect two nodes with an edge, get the pointer to the edge
func (g *Graph) AddEdge(a, b *Node) *Edge {
	aToB := &Edge{
		Tail: a,
		Head: b,
	}
	g.Edges.PushBack(aToB)
	a.Edges.PushBack(aToB)
	b.Edges.PushBack(aToB)
	return aToB
}

func NewNode() Node {
	return Node {
		Position: rl.Vector2Zero(),
		Contents: "",
		Edges: list.New(),
	}
}

func (n *Node) IsConnectedTo(b *Node) bool {
	for edgeIt := n.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*Edge)
		if edge.Head == b {
			return true
		}
	}
	return false

}
