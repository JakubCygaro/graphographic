package algorithm

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// import "fmt"

type elem[T any] struct {
	key int32
	val T
}

type MinHeap[T any] = []elem[T]

func Len[T any](heap *MinHeap[T]) int {
	return len(*heap)
}

func GetMin[T any](heap *MinHeap[T]) (int32, *T) {
	var ret *T = nil
	if Len(heap) > 0 {
		ret = &(*heap)[0].val
	}
	return (*heap)[0].key, ret
}

func Pop[T any](heap *MinHeap[T]) {
	if Len(heap) > 0 {
		swap(heap, 0, Len(heap) - 1)
		*heap = (*heap)[:Len(heap) - 1]
	}
	siftDown(heap, 0)
}

func Insert[T any](heap *MinHeap[T], k int32, v T) {
	*heap = append(*heap, elem[T]{ key: k, val: v })
	siftUp(heap, Len(heap) - 1)
	verifyHeapProp(heap)
}

func Search[T comparable](heap *MinHeap[T], v T) (int32, bool) {
	found := -1
	for i, elem := range *heap {
		if elem.val == v {
			found = i
			break
		}
	}
	if found != -1 {
		return (*heap)[found].key, true
	}
	return int32(0), false
}

func Delete[T comparable](heap *MinHeap[T], v T) {
	found := -1
	for i, elem := range *heap {
		if elem.val == v {
			found = i
			break
		}
	}
	if found != -1 {
		//swap found with last element
		swap(heap, found, Len(heap) - 1)
		// remove found from heap
		*heap = (*heap)[:Len(heap)-1]
		// sift swapped element up or down to maintain heap property
		if found != Len(heap) {
			found = siftUp(heap, found)
			siftDown(heap, found)
		}
	}
	verifyHeapProp(heap)
}

func parent(idx int) int {
	if idx == 0 {
		return -1
	}
	return (idx - 1) / 2
}
func leftChild[T any](heap *MinHeap[T], idx int) int {
	ret := idx*2 + 1
	if ret >= len(*heap) {
		return -1
	} else {
		return ret
	}
}
func rightChild[T any](heap *MinHeap[T], idx int) int {
	ret := idx*2 + 2
	if ret >= len(*heap) {
		return -1
	} else {
		return ret
	}
}

func siftDown[T any](heap *MinHeap[T], elem int) int {
	if elem == -1 {
		return -1
	}
	for {
		left, right := leftChild(heap, elem), rightChild(heap, elem)
		var candidate int = -1
		// if elem has a left child, it could have a right one, but without a left child there is no right child
		candidate = left
		if right != -1 && (*heap)[right].key <= (*heap)[left].key {
			candidate = right
		}
		if candidate != -1 && (*heap)[candidate].key < (*heap)[elem].key {
			swap(heap, elem, candidate)
			elem = candidate
		} else {
			return elem
		}
	}
}

func siftUp[T any](heap *MinHeap[T], elem int) int {
	if elem == -1 {
		return -1
	}
	for {
		parent := parent(elem)
		if parent != -1 && (*heap)[parent].key > (*heap)[elem].key {
			swap(heap, parent, elem)
			elem = parent
		} else {
			return elem
		}
	}
}

func verifyHeapProp[T any](heap* MinHeap[T]){
	for i, e := range *heap {
		p := parent(i)
		l, r := leftChild(heap, i), rightChild(heap, i)
		if p != -1 && (*heap)[p].key > e.key {
			rl.TraceLog(rl.LogError, "PARENT HEAP PROPERTY VIOLATED IDX: %d", i)
			printKeys(heap)
		}
		if l != -1 && (*heap)[l].key < e.key {
			rl.TraceLog(rl.LogError, "LEFT HEAP PROPERTY VIOLATED IDX: %d", i)
			printKeys(heap)
		}
		if r != -1 && (*heap)[r].key < e.key {
			rl.TraceLog(rl.LogError, "RIGHT HEAP PROPERTY VIOLATED IDX: %d", i)
			printKeys(heap)
		}
	}
}

func swap[T any](heap* MinHeap[T], a, b int){
	(*heap)[a], (*heap)[b] = (*heap)[b], (*heap)[a]
}
func printKeys[T any](heap* MinHeap[T]){
	for _, e := range *heap {
		fmt.Print(e.key, " ")
	}
	fmt.Println()
}
