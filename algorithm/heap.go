package algorithm

// import "fmt"

type elem[T any] struct {
	key int32
	val T
}

type MinHeap[T any] = []elem[T]

func Len[T any](heap *MinHeap[T]) int {
	return len(*heap)
}

func GetMin[T any](heap *MinHeap[T]) *T {
	var ret *T = nil
	if Len(heap) > 0 {
		ret = &(*heap)[0].val
	}
	return ret
}

func Pop[T any](heap *MinHeap[T]) {
	if Len(heap) > 0 {
		*heap = (*heap)[1:]
	}
}

func Insert[T any](heap *MinHeap[T], k int32, v T) {
	*heap = append(*heap, elem[T]{ key: k, val: v })
	siftUp(heap, Len(heap) - 1)
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
		(*heap)[found], (*heap)[Len(heap)-1] = (*heap)[Len(heap)-1], (*heap)[found]
		// remove found from heap
		*heap = (*heap)[:Len(heap)-1]
		// sift swapped element up or down to maintain heap property
		if found != Len(heap) {
			found = siftUp(heap, found)
			siftDown(heap, found)
		}
	}
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
		if left != -1 && right == -1 {
			candidate = left
		} else if right != -1 && (*heap)[right].key < (*heap)[left].key {
			candidate = right
		}
		if candidate != -1 && (*heap)[candidate].key < (*heap)[elem].key {
			(*heap)[elem], (*heap)[candidate] = (*heap)[candidate], (*heap)[elem]
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
			(*heap)[parent], (*heap)[elem] = (*heap)[elem], (*heap)[parent]
			elem = parent
		} else {
			return elem
		}
	}
}
