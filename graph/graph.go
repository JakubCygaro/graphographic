package graph

import (
	"container/list"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type AlgoData struct {
	Explored bool
	Highlighted bool
	// extra data that could be assigned and used by an algorithm
	Custom *any
}

type Node struct {
	Position rl.Vector2
	Content string
	Edges *list.List
	// saved from draw pass
	Radius float32
	Data AlgoData
}

type Edge struct {
	Tail *Node
	Head *Node
	Cost int32
	// set after each draw pass so it does not have to be recalculated
	StartPos, EndPos rl.Vector2
	Data AlgoData
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
	return aToB
}

func (node *Node) removeEdgeToNode(other *Node) {
	edgesToRemove := make([]*list.Element, 0)
	for edgeIt := node.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*Edge)
		if edge.Tail == other || edge.Head == other {
			edgesToRemove = append(edgesToRemove, edgeIt)
		}
	}
	for _, e := range edgesToRemove {
		node.Edges.Remove(e)
	}
}

func (node *Node) removeEdge(e *Edge) {
	for edgeIt := node.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*Edge)
		if edge == e {
			node.Edges.Remove(edgeIt)
		}
	}
}

func (g *Graph) RemoveNode(n *Node) {
	edgesToRemove := make([]*list.Element, 0)
	for edgeIt := g.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*Edge)
		if edge.Tail == n || edge.Head == n {
			edgesToRemove = append(edgesToRemove, edgeIt)
			edge.Tail.removeEdgeToNode(n)
			edge.Head.removeEdgeToNode(n)
		}
	}
	for _, e := range edgesToRemove {
		g.Edges.Remove(e)
	}

	for nodeIt := g.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		node := nodeIt.Value.(*Node)
		if node == n {
			g.Nodes.Remove(nodeIt)
			break;
		}
	}
}
func (g *Graph) AddNode(n Node) *Node {
	nPtr := &n
	g.Nodes.PushBack(nPtr)
	return nPtr
}
func (g *Graph) RemoveEdge(e *Edge) {
	for edgeIt := g.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		edge := edgeIt.Value.(*Edge)
		if edge == e {
			edge.Head.removeEdge(edge)
			edge.Tail.removeEdge(edge)
			g.Edges.Remove(edgeIt)
			break
		}
	}
}

func NewNode() Node {
	return Node {
		Position: rl.Vector2Zero(),
		Content: "",
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
