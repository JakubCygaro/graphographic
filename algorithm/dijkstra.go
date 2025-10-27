package algorithm

import (
	"fmt"
	"graphographic/graph"
	"math"
)

type data struct {
	Cost int32
	Prev *graph.Node
}

type Dijkstra struct {
	heap  MinHeap[*graph.Node]
	start *graph.Node
	end   *graph.Node
	last  *graph.Node
}

func (algo *Dijkstra) Init() {
	algo.heap = make(MinHeap[*graph.Node], 0)
	algo.start = nil
	algo.end = nil
	algo.last = nil
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
	d := new(data)
	d.Prev = nil
	d.Cost = 0
	algo.start.Data.Custom = d
	algo.start.Data.Explored = true
	algo.last = algo.start
	for nodeIt := g.Nodes.Front(); nodeIt != nil; nodeIt = nodeIt.Next() {
		n := nodeIt.Value.(*graph.Node)
		d := new(data)
		d.Cost = math.MaxInt32
		d.Prev = nil
		n.Data.Custom = d
		n.Data.Explored = false
		Insert(&algo.heap, d.Cost, n)
	}
	for edgeIt := algo.start.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
		e := edgeIt.Value.(*graph.Edge)
		if e.Tail == algo.start {
			d := e.Head.Data.Custom.(*data)
			d.Cost = e.Cost
			d.Prev = algo.start
			Delete(&algo.heap, e.Head)
			Insert(&algo.heap, d.Cost, e.Head)
		}
	}

	return nil
}

func (algo *Dijkstra) Update() bool {
	if Len(&algo.heap) > 0 {
		next := *GetMin(&algo.heap)
		Pop(&algo.heap)
		d := next.Data.Custom.(*data)
		next.Data.Explored = true
		for edgeIt := next.Edges.Front(); edgeIt != nil; edgeIt = edgeIt.Next() {
			edge := edgeIt.Value.(*graph.Edge)
			if !edge.Head.Data.Explored {
				currentCost, _ := Search(&algo.heap, edge.Head)
				if currentCost > d.Cost+edge.Cost {
					Delete(&algo.heap, edge.Head)
					d2 := edge.Head.Data.Custom.(*data)
					d2.Cost = d.Cost + edge.Cost
					Insert(&algo.heap, d2.Cost, edge.Head)
				}
			}
		}
		d.Prev = algo.last
		algo.last = next
		if next == algo.end {
			return false
		}
		return true
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
