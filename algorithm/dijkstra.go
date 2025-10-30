package algorithm

import (
	"fmt"
	"graphographic/graph"
	"math"
)

type data struct {
	Len int32
	Prev *graph.Node
	InPath bool
}

type Dijkstra struct {
	heap  MinHeap[*graph.Node]
	start *graph.Node
	end   *graph.Node
	prev  *graph.Node
	graph *graph.Graph
}

func (algo *Dijkstra) Init() {
	algo.heap = make(MinHeap[*graph.Node], 0)
	algo.start = nil
	algo.end = nil
	algo.prev = nil
}

func (algo *Dijkstra) GetName() string {
	return "Dijkstra"
}
func (algo *Dijkstra) Start(g *graph.Graph) error {
	if algo.start == nil || algo.end == nil {
		algo.start = nil
		algo.end = nil
		return fmt.Errorf("Start or End node not selected")
	}
	algo.graph = g
	startData := new(data)
	startData.Prev = nil
	startData.Len = 0
	startData.InPath = false
	algo.start.Data.Custom = startData

	algo.start.Data.Highlighted = false
	algo.start.Data.Explored = true

	Insert(&algo.heap, 0, algo.start)

	for nodeIt := g.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		n := nodeIt.Value.(*graph.Node)
		if n == algo.start {
			continue
		}
		nodeData := new(data)
		nodeData.Len = math.MaxInt32
		nodeData.Prev = nil
		nodeData.InPath = false
		n.Data.Custom = nodeData

		n.Data.Explored = false
		Insert(&algo.heap, nodeData.Len, n)
	}
	return nil
}

func (algo *Dijkstra) Update() bool {
	if Len(&algo.heap) > 0 {
		k, n := GetMin(&algo.heap)
		next := *n
		Pop(&algo.heap)
		nextNodeData := next.Data.Custom.(*data)
		next.Data.Explored = true
		nextNodeData.InPath = true
		nextNodeData.Len = k
		next.Data.Tag = fmt.Sprintf("%d", nextNodeData.Len)
		for edgeIt := next.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
			edge := edgeIt.Value.(*graph.Edge)
			headData := edge.Head.Data.Custom.(*data)
			if !headData.InPath && nextNodeData.Len != math.MaxInt32{
				currentCost, _ := Search(&algo.heap, edge.Head)
				if currentCost > nextNodeData.Len+edge.Cost {
					Delete(&algo.heap, edge.Head)
					d2 := edge.Head.Data.Custom.(*data)
					d2.Len = nextNodeData.Len + edge.Cost
					headData.Prev = next
					Insert(&algo.heap, d2.Len, edge.Head)
				}
			}
		}
		// nextNodeData.Prev = algo.prev
		// if nextNodeData.Len != math.MaxInt32 {
		// 	algo.prev = next
		// }
		if next != algo.end {
			return true
		}
		algo.prev = next
	}
	if algo.prev == algo.end {
		for nodeIt := algo.graph.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
			n := nodeIt.Value.(*graph.Node)
			n.Data.Explored = false
		}
		for prev := algo.prev; prev != nil; {
			fmt.Println(prev.Content)
			prev.Data.Explored = true
			d := prev.Data.Custom.(*data)
			prev = d.Prev
		}
	}
	return false
}

func (algo *Dijkstra) NodeSelected(node *graph.Node) {
	if algo.start == nil {
		algo.start = node
		algo.start.Data.Highlighted = true
	} else if algo.end == nil {
		algo.end = node
		algo.end.Data.Highlighted = true
	}
}
func (algo *Dijkstra) UndoSelect() {
	if algo.end != nil {
		algo.end.Data.Highlighted = false
	} else if algo.start != nil {
		algo.start.Data.Highlighted = false
	}
}
