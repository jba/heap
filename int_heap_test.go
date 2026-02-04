// Non-generic int heap implementation for benchmarking comparison.
// Uses direct < comparisons instead of cmp.Compare.

package heap

import (
	"cmp"
	"slices"
	"testing"
)

// intHeap is a min-heap for int values.
type intHeap struct {
	values []int
}

// newIntHeap creates a new min-heap for int values.
func newIntHeap() *intHeap {
	return &intHeap{}
}

// Insert adds an element to the heap.
func (h *intHeap) Insert(value int) {
	h.values = append(h.values, value)
	h.up(len(h.values) - 1)
}

func (h *intHeap) Init(s []int) {
	if len(h.values) != 0 {
		panic("non-empty")
	}
	h.values = s
	h.build()
}

// Min returns the minimum element in the heap without removing it.
func (h *intHeap) Min() int {
	return h.values[0]
}

func (h *intHeap) TakeMin() int {
	if len(h.values) == 0 {
		panic("heap: TakeMin called on empty heap")
	}
	min := h.values[0]
	h.delete(0)
	return min
}

func (h *intHeap) delete(i int) {
	n := len(h.values) - 1
	if n != i {
		h.values[i], h.values[n] = h.values[n], h.values[i]
	}
	h.values[n] = 0 // allow GC
	h.values = h.values[:n]
	if n != i && !h.down(i) {
		h.up(i)
	}
}

// ChangeMin replaces the minimum value in the heap with the given value.
func (h *intHeap) ChangeMin(v int) {
	h.values[0] = v
	h.down(0)
}

func (h *intHeap) build() {
	n := len(h.values)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i)
	}
}

// Len returns the number of elements in the heap.
func (h *intHeap) Len() int {
	return len(h.values)
}

// up moves the element at index i up the heap until the heap invariant is restored.
func (h *intHeap) up(i int) {
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
func (h *intHeap) down(i int) bool {
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

func TestIntHeapSort(t *testing.T) {
	h := newIntHeap()
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

func TestIntHeapTopK(t *testing.T) {
	data := []int{7, 2, 9, 1, 5, 8, 3, 6, 4, 10}
	const k = 3

	h := newIntHeap()
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
