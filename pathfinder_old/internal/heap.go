package internal

import (
	"log"
)

type NodesMinHeap struct {
	arr      []*Node
	heapSize int
}

func NewHeap(cap int) *NodesMinHeap {
	heap := NodesMinHeap{
		arr:      make([]*Node, cap),
		heapSize: 0,
	}

	return &heap
}

func (h *NodesMinHeap) parent(i int) int {
	return (i - 1) / 2
}

func (h *NodesMinHeap) left(i int) int {
	return 2*i + 1
}

func (h *NodesMinHeap) right(i int) int {
	return 2*i + 2
}

func (h *NodesMinHeap) swap(i1, i2 int) {
	temp := h.arr[i1]
	h.arr[i1] = h.arr[i2]
	h.arr[i2] = temp
}

func (h *NodesMinHeap) minimize(i int) {
	l := h.left(i)
	r := h.right(i)
	smallest := i

	if l < h.heapSize && h.arr[l].GetScore() < h.arr[i].GetScore() {
		smallest = l
	}

	if r < h.heapSize && h.arr[r].GetScore() < h.arr[smallest].GetScore() {
		smallest = r
	}

	if smallest != i {
		h.swap(i, smallest)
		h.minimize(smallest)
	}
}

func (h *NodesMinHeap) Empty() bool {
	return h.heapSize == 0
}

func (h *NodesMinHeap) Clear() {
	h.heapSize = 0
}

func (h *NodesMinHeap) Insert(node *Node) {
	if h.heapSize == cap(h.arr) {
		log.Println("warning: heap is full, reallocation")
		h.arr = append(h.arr, node)
		newArr := make([]*Node, cap(h.arr)-len(h.arr))
		h.arr = append(h.arr, newArr...)
	}

	h.heapSize++
	i := h.heapSize - 1
	h.arr[i] = node

	for i != 0 && h.arr[h.parent(i)].GetScore() > h.arr[i].GetScore() {
		h.swap(i, h.parent(i))
		i = h.parent(i)
	}
}

func (h *NodesMinHeap) GetMin() *Node {
	if h.heapSize == 0 {
		return nil
	}

	if h.heapSize == 1 {
		h.heapSize--
		return h.arr[0]
	}

	root := h.arr[0]
	h.arr[0] = h.arr[h.heapSize-1]
	h.heapSize--
	h.minimize(0)

	return root
}
