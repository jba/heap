// Heap[T cmp.Ordered] implementation for benchmarking comparison.
// Simplified version without mover interface.

package heap

import (
	"cmp"
	"slices"
	"testing"
)

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

func (h *orderedHeap[T]) Init(s []T) {
	if len(h.values) != 0 {
		panic("non-empty")
	}
	h.values = s
	h.build()
}

// Min returns the minimum element in the heap without removing it.
func (h *orderedHeap[T]) Min() T {
	return h.values[0]
}

func (h *orderedHeap[T]) TakeMin() T {
	if len(h.values) == 0 {
		panic("heap: TakeMin called on empty heap")
	}
	min := h.values[0]
	h.delete(0)
	return min
}

func (h *orderedHeap[T]) delete(i int) {
	n := len(h.values) - 1
	if n != i {
		h.values[i], h.values[n] = h.values[n], h.values[i]
	}
	var zero T
	h.values[n] = zero // allow GC
	h.values = h.values[:n]
	if n != i && !h.down(i) {
		h.up(i)
	}
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
func (h *orderedHeap[T]) down(i int) bool {
	data := h.values
	n := len(data)
	i0 := i
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
	return i > i0
}

func TestOrderedHeapSort(t *testing.T) {
	h := newOrderedHeap[int]()
	h.Init([]int{5, 2, 8, 1, 9, 3, 7, 2, 7})

	var got []int
	for h.Len() > 0 {
		got = append(got, h.TakeMin())
	}

	want := []int{1, 2, 2, 3, 5, 7, 7, 8, 9}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestOrderedHeapTopK(t *testing.T) {
	data := []int{7, 2, 9, 1, 5, 8, 3, 6, 4, 10}
	const k = 3

	h := newOrderedHeap[int]()
	h.Init(slices.Clone(data[:k]))

	for _, v := range data[k:] {
		if v > h.Min() {
			h.ChangeMin(v)
		}
	}

	var got []int
	for h.Len() > 0 {
		got = append(got, h.TakeMin())
	}

	want := []int{8, 9, 10}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
