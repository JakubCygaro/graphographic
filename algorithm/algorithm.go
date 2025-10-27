package algorithm

import (
	"graphographic/graph"
)

type Algorithm interface {
	Init();
	NodeSelected(node *graph.Node)
	UndoSelect()
	Start(*graph.Graph) error;
	Update() bool;
	GetName() string;
}

