package algorithm
import (
	"fmt"
	"graphographic/graph"
)

type DFS struct {
	start *graph.Node
	queue []*graph.Node
}
func (algo *DFS) Init() {
	algo.start = nil
}
func (algo *DFS) GetName() string {
	return "DFS"
}
func (algo *DFS) Start(g *graph.Graph) error {
	if algo.start == nil {
		return fmt.Errorf("Starting node was not selected")
	}
	algo.queue = make([]*graph.Node, 0)
	algo.queue = append(algo.queue, algo.start)
	return nil
}

func (algo *DFS) Update() bool {
	var next *graph.Node = nil
	if len(algo.queue) > 0 {
		next, algo.queue = algo.queue[0], algo.queue[1:]
		algo.addNodesToQueue(next)
		return true
	} else {
		return false
	}
}

func (algo *DFS) addNodesToQueue(node *graph.Node){
	node.Data.Explored = true
	for edgeIt := node.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		e := edgeIt.Value.(*graph.Edge)
		if e.Tail == node && !e.Head.Data.Explored {
			e.Data.Explored = true
			algo.queue = append(algo.queue, e.Head)
		}
	}
}

func (algo *DFS) NodeSelected(node *graph.Node) {
	if algo.start != nil {
		algo.start.Data.Highlighted = false
	}
	if algo.start = node; algo.start != nil {
		algo.start.Data.Highlighted = true
	}
}
func (algo *DFS) UndoSelect(){
	if algo.start != nil {
		algo.start.Data.Highlighted = false
	}
}
