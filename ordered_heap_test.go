// Heap[T cmp.Ordered] implementation for benchmarking comparison.
// Simplified version without mover interface.

package heap

import "cmp"

// orderedHeap is a min-heap for ordered types.
type orderedHeap[T cmp.Ordered] struct {
	values []T
}

// newOrderedHeap creates a new min-heap for ordered types.
func newOrderedHeap[T cmp.Ordered]() *orderedHeap[T] {
	return &orderedHeap[T]{}
}

// Insert adds an element to the heap.
func (h *orderedHeap[T]) Insert(value T) {
	h.values = append(h.values, value)
	h.up(len(h.values) - 1)
}

// InsertSlice adds all elements of s to the heap, then heapifies.
func (h *orderedHeap[T]) InsertSlice(s []T) {
	if h.values == nil {
		h.values = s
	} else {
		h.values = append(h.values, s...)
	}
	h.build()
}

// Min returns the minimum element in the heap without removing it.
func (h *orderedHeap[T]) Min() T {
	return h.values[0]
}

// ChangeMin replaces the minimum value in the heap with the given value.
func (h *orderedHeap[T]) ChangeMin(v T) {
	h.values[0] = v
	h.down(0)
}

func (h *orderedHeap[T]) build() {
	n := len(h.values)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i)
	}
}

// Len returns the number of elements in the heap.
func (h *orderedHeap[T]) Len() int {
	return len(h.values)
}

// up moves the element at index i up the heap until the heap invariant is restored.
func (h *orderedHeap[T]) up(i int) {
	for i > 0 {
		p := (i - 1) / 2 // parent
		if cmp.Compare(h.values[i], h.values[p]) >= 0 {
			break
		}
		h.values[p], h.values[i] = h.values[i], h.values[p]
		i = p
	}
}

// down moves the element at index i down the heap until the heap invariant is restored.
func (h *orderedHeap[T]) down(i int) {
	data := h.values
	n := len(data)
	for {
		lc := 2*i + 1
		if lc >= n {
			break
		}
		child := lc // left child
		if rc := lc + 1; rc < n && cmp.Compare(data[rc], data[lc]) < 0 {
			child = rc // right child is smaller
		}
		if cmp.Compare(data[child], data[i]) >= 0 {
			break
		}
		data[i], data[child] = data[child], data[i]
		i = child
	}
}
