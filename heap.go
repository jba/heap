// Package heap provides a min-heap data structure.
package heap

import "iter"

// A Heap is a binary min-heap.
type Heap[T any] struct {
	values   []T
	setIndex func(T, int)
	compare  func(T, T) int
}

// New creates a new [Heap] with the given comparison function.
// The comparison function should return:
//   - a negative value if a < b
//   - zero if a == b
//   - a positive value if a > b.
func New[T any](compare func(T, T) int) *Heap[T] {
	return &Heap[T]{compare: compare}
}

// SetIndexFunc sets a function that sets the index of a heap element.
// It should be called before any other function on the heap.
// A Heap with an index function supports the Delete and Changed methods.
// In a Heap with an index function, all the elements must be distinct.
func (h *Heap[T]) SetIndexFunc(f func(T, int)) {
	h.setIndex = f
}

// Insert adds an element to the heap.
func (h *Heap[T]) Insert(value T) {
	h.values = append(h.values, value)
	if h.setIndex != nil {
		h.setIndex(value, len(h.values)-1)
	}
	h.up(len(h.values) - 1)
}

// InsertSlice adds all elements of s to the heap, then heapifies.
// The caller must not subsequently modify s.
func (h *Heap[T]) InsertSlice(s []T) {
	start := len(h.values)
	if start == 0 {
		h.values = s
	} else {
		h.values = append(h.values, s...)
	}
	if h.setIndex != nil {
		for i, v := range s {
			h.setIndex(v, start+i)
		}
	}
	// heapify
	for i := len(h.values)/2 - 1; i >= 0; i-- {
		h.down(i)
	}
}

// Min returns the minimum element in the heap without removing it.
// It panics if the heap is empty.
func (h *Heap[T]) Min() T {
	if len(h.values) == 0 {
		panic("heap: Min called on empty heap")
	}
	return h.values[0]
}

// TakeMin removes and returns the minimum element from the heap.
// It panics if the heap is empty.
func (h *Heap[T]) TakeMin() T {
	if len(h.values) == 0 {
		panic("heap: TakeMin called on empty heap")
	}
	min := h.values[0]
	h.delete(0)
	return min
}

// Clear removes all elements from the heap.
func (h *Heap[T]) Clear() {
	if h.setIndex != nil {
		for _, v := range h.values {
			h.setIndex(v, -1)
		}
	}
	var zero T
	for i := range h.values {
		h.values[i] = zero // allow GC
	}
	h.values = h.values[:0]
}

// Len returns the number of elements in the heap.
func (h *Heap[T]) Len() int {
	return len(h.values)
}

// All returns an iterator over all elements in the heap
// in unspecified order.
func (h *Heap[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range h.values {
			if !yield(v) {
				return
			}
		}
	}
}

// Drain removes and returns the heap elements in sorted order,
// from smallest to largest.
//
// The result is undefined if the heap is changed during iteration.
func (h *Heap[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for len(h.values) > 0 {
			if !yield(h.TakeMin()) {
				return
			}
		}
	}
}

// Delete removes the element at index i from the heap.
// The only reasonable values for i are 0, for the minimum element (but
// see [Heap.TakeMin]),
// or an index maintained by an index function (see [Heap.SetIndexFunc]).
// If i is out of range, or it is non-zero and there is no index function,
// Delete panics.
func (h *Heap[T]) Delete(i int) {
	if i < 0 || i >= len(h.values) {
		panic("heap: Delete: index out of range")
	}
	if i != 0 && h.setIndex == nil {
		panic("heap: Delete called with non-zero index and no index function")
	}
	h.delete(i)
}

func (h *Heap[T]) delete(i int) {
	if h.setIndex != nil {
		h.setIndex(h.values[i], -1)
	}
	n := len(h.values) - 1
	if n != i {
		h.swap(i, n)
	}
	var zero T
	h.values[n] = zero // allow GC
	h.values = h.values[:n]
	if n != i && !h.down(i) {
		h.up(i)
	}
}

// Changed restores the heap property after the element at index i has
// been modified. The only reasonable values for i are 0, for the minimum
// element (but see [Heap.ChangeMin] for an alternative) or an index maintained
// by an index function (see [Heap.SetIndexFunc]). If i is out of range,
// or it is non-zero and there is no index function, Changed panics.
func (h *Heap[T]) Changed(i int) {
	if i < 0 || i >= len(h.values) {
		panic("heap: Changed: index out of range")
	}
	if i != 0 && h.setIndex == nil {
		panic("heap: Changed called with non-zero index and no index function")
	}
	if !h.down(i) {
		h.up(i)
	}
}

// ChangeMin replaces the minimum value in the heap with the given value.
// It panics if the heap is empty.
func (h *Heap[T]) ChangeMin(v T) {
	if len(h.values) == 0 {
		panic("heap: ChangeMin called on empty heap")
	}
	h.values[0] = v
	h.down(0)
}

// up moves the element at index i up the heap until the heap property
// is restored. func (h *Heap[T]) up(i int) {
// 	for i > 0 {
// 		p := (i - 1) / 2 // ], h.values[p]) >= 0 {
// 			break
		}
		h.swap(p, i)
		i = p
	}
}

// down moves the element at index i down the heap until the heap property
// is restored. Returns true if the element moved. func (h *Heap[T]) down(i
// int) bool {
// 	n := le], h.values[lc]) < 0 {
			child = rc // right child is smaller
		}
		if h.compare(h.values[child], h.values[i]) >= 0 {
			break
		}
		h.swap(i, child)
		i = child
	}
	return i > i0
}

func (h *Heap[T]) swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
	if h.setIndex != nil {
		h.setIndex(h.values[i], i)
		h.setIndex(h.values[j], j)
	}
}
