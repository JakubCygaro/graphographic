package algorithm

import (
	"graphographic/graph"
)

type Algorithm interface {
	Init();
	NodeSelected(node *graph.Node)
	Start(*graph.Graph) error;
	Update() bool;
	GetName() string;
}

