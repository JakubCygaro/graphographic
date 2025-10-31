package algorithm

import (
	"fmt"
	"graphographic/graph"
)

type BFS struct {
	start *graph.Node
	stack []*graph.Node
}

func (algo *BFS) Init() {
	algo.start = nil
}
func (algo *BFS) GetName() string {
	return "BFS"
}
func (algo *BFS) Start(g *graph.Graph) error {
	if algo.start == nil {
		return fmt.Errorf("Starting node was not selected")
	}
	algo.stack = make([]*graph.Node, 0)
	algo.stack = append(algo.stack, algo.start)
	return nil
}

func (algo *BFS) Update() bool {
	var next *graph.Node = nil
	if len(algo.stack) > 0 {
		next, algo.stack = algo.stack[len(algo.stack)-1], algo.stack[:len(algo.stack)-1]
		algo.addNodesToStack(next)
		return true
	} else {
		return false
	}
}

func (algo *BFS) addNodesToStack(node *graph.Node) {
	node.Data.Explored = true
	for edgeIt := node.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		e := edgeIt.Value.(*graph.Edge)
		if e.Tail == node && !e.Head.Data.Explored {
			e.Data.Explored = true
			algo.stack = append(algo.stack, e.Head)
		}
	}
}

func (algo *BFS) NodeSelected(node *graph.Node) {
	if algo.start != nil {
		algo.start.Data.Highlighted = false
	}
	if algo.start = node; algo.start != nil {
		algo.start.Data.Highlighted = true
	}
}

func (algo *BFS) UndoSelect(){
	if algo.start != nil {
		algo.start.Data.Highlighted = false
	}
}
