// Package heap provides min-heap data structures.
package heap

import (
	"iter"
)

// Heap is a min-heap for any type with a custom comparison function.
type Heap[T any] struct {
	values    []T
	indexes   []*int // from calls to indexFunc; updated in swap
	indexFunc func(T) *int
	compare   func(T, T) int
}

// New creates a new min-heap with a custom comparison function.
// The comparison function should return a negative value if a < b,
// zero if a == b, and a positive value if a > b.
func New[T any](compare func(T, T) int) *Heap[T] {
	return &Heap[T]{compare: compare}
}

// SetIndexFunc sets a function that returns a pointer to an index field
// in the heap element.
func (h *Heap[T]) SetIndexFunc(f func(T) *int) {
	h.indexFunc = f
	h.indexes = make([]*int, len(h.values))
}

// Insert adds an element to the heap.
func (h *Heap[T]) Insert(value T) {
	h.values = append(h.values, value)
	if h.indexes != nil {
		p := h.indexFunc(value)
		*p = len(h.indexes)
		h.indexes = append(h.indexes, p)
	}
	h.up(len(h.values) - 1)
}

// InsertSlice adds all elements of s to the heap, then heapifies.
// The caller must not subsequently modify s.
func (h *Heap[T]) InsertSlice(s []T) {
	if h.values == nil {
		h.values = s
	} else {
		h.values = append(h.values, s...)
	}
	if h.indexFunc != nil {
		start := len(h.indexes)
		for i, v := range s {
			p := h.indexFunc(v)
			*p = start + i
			h.indexes = append(h.indexes, p)
		}
	}
	h.build()
}

// Min returns the minimum element in the heap without removing it.
// It panics if the heap is empty.
func (h *Heap[T]) Min() T {
	return h.min()
}

func (h *Heap[T]) min() T {
	if len(h.values) == 0 {
		panic("heap: Min called on empty heap")
	}
	return h.values[0]
}

// TakeMin removes and returns the minimum element from the heap.
// It panics if the heap is empty.
func (h *Heap[T]) TakeMin() T {
	return h.takeMin()
}

func (h *Heap[T]) takeMin() T {
	if len(h.values) == 0 {
		panic("heap: TakeMin called on empty heap")
	}
	min := h.values[0]
	h.deleteAt(0)
	return min
}

// ChangeMin replaces the minimum value in the heap with the given value.
// It panics if the heap is empty.
func (h *Heap[T]) ChangeMin(v T) {
	h.changeMin(v)
}

func (h *Heap[T]) changeMin(v T) {
	if len(h.values) == 0 {
		panic("heap: ChangeMin called on empty heap")
	}
	h.values[0] = v
	h.down(0)
}

func (h *Heap[T]) build() {
	n := len(h.values)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i)
	}
}

// Clear removes all elements from the heap.
func (h *Heap[T]) Clear() {
	h.clear()
}

func (h *Heap[T]) clear() {
	var zero T
	for i := range h.values {
		h.values[i] = zero // allow GC
	}
	h.values = h.values[:0]
	if h.indexes != nil {
		for i := range h.indexes {
			*h.indexes[i] = -1
			h.indexes[i] = nil // allow GC
		}
		h.indexes = h.indexes[:0]
	}
}

// Len returns the number of elements in the heap.
func (h *Heap[T]) Len() int {
	return len(h.values)
}

// All returns an iterator over all elements in the heap
// in unspecified order.
func (h *Heap[T]) All() iter.Seq[T] {
	return h.all()
}

func (h *Heap[T]) all() iter.Seq[T] {
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
	return h.drain()
}

func (h *Heap[T]) drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for len(h.values) > 0 {
			if !yield(h.takeMin()) {
				return
			}
		}
	}
}

// Delete removes the item from the heap.
// If the item has already been deleted or the heap has been cleared,
// Delete does nothing.
//
// The Heap must have an index function.
func (h *Heap[T]) Delete(v T) {
	h.delete(v)
}

func (h *Heap[T]) delete(v T) {
	if h.indexFunc == nil {
		panic("heap: Delete: SetIndexFunc was not called")
	}
	pi := h.indexFunc(v)
	if *pi < 0 || *pi >= len(h.values) {
		return
	}
	// Sanity check: the entry at v's index should point to the same place.
	if h.indexes[*pi] != pi {
		panic("heap: Delete: index pointer mismatch")
	}
	h.deleteAt(*pi)
}

func (h *Heap[T]) deleteAt(i int) {
	if h.indexes != nil {
		*h.indexes[i] = -1
	}
	n := len(h.values) - 1
	if n != i {
		h.swap(i, n)
	}
	var zero T
	h.values[n] = zero // allow GC
	h.values = h.values[:n]
	if h.indexes != nil {
		h.indexes[n] = nil // allow GC
		h.indexes = h.indexes[:n]
	}
	if n != i && !h.down(i) {
		h.up(i)
	}
}

// Changed restores the heap invariant after the item's value has been changed.
// Call this method after modifying the value of the item.
// If the item has been deleted or the heap has been cleared, Changed does nothing.
//
// The Heap must have an index function.
func (h *Heap[T]) Changed(v T) {
	h.changed(v)
}

func (h *Heap[T]) changed(v T) {
	if h.indexFunc == nil {
		panic("heap: Changed: no index function")
	}
	ip := h.indexFunc(v)
	if *ip < 0 || *ip >= len(h.values) {
		return
	}
	// Sanity check: the entry at v's index should point to the same place.
	if h.indexes[*ip] != ip {
		panic("heap: Changed: index pointer mismatch")
	}
	if !h.down(*ip) {
		h.up(*ip)
	}
}

// up moves the element at index i up the heap until the heap invariant is restored.
func (h *Heap[T]) up(i int) {
	for i > 0 {
		p := (i - 1) / 2 // parent
		if h.compare(h.values[i], h.values[p]) >= 0 {
			break
		}
		h.swap(p, i)
		i = p
	}
}

// down moves the element at index i down the heap until the heap invariant is restored.
// Returns true if the element moved.
func (h *Heap[T]) down(i int) bool {
	n := len(h.values)
	i0 := i
	for {
		lc := 2*i + 1
		if lc >= n {
			break
		}
		child := lc // left child
		if rc := lc + 1; rc < n && h.compare(h.values[rc], h.values[lc]) < 0 {
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
	if h.indexes != nil {
		h.indexes[i], h.indexes[j] = h.indexes[j], h.indexes[i]
		*h.indexes[i] = i
		*h.indexes[j] = j
	}
}
