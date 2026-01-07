// Package heap provides min-heap data structures.
package heap

import (
	"cmp"
	"iter"
)

// A Heap is a binary min-heap.
type Heap[T any] interface {
	// Insert adds an element to the heap, preserving the heap property.
	Insert(value T)
	// InsertSlice adds all elements of s to the heap, then re-establishes
	// the heap property.
	// The caller must not subsequently modify s.
	InsertSlice(s []T)
	// Min returns the minimum element in the heap without removing it.
	// It panics if the heap is empty.
	Min() T
	// TakeMin removes and returns the minimum element from the heap.
	// It panics if the heap is empty.
	TakeMin() T
	// ChangeMin replaces the minimum value in the heap with the given value.
	// It panics if the heap is empty.
	ChangeMin(v T)
	// Len returns the number of elements in the heap.
	Len() int
	// All returns an iterator over all elements in the heap
	// in unspecified order.
	All() iter.Seq[T]
	// Drain removes and returns the heap elements in sorted order,
	// from smallest to largest.
	//
	// The result is undefined if the heap is changed during iteration.
	Drain() iter.Seq[T]
	// Clear removes all elements from the heap.
	Clear()

	// SetIndexFunc sets a function that returns a pointer index for
	// the given heap element. The index function must not return nil.
	//
	// SetIndexFunc enables the use of the Delete and Changed methods.
	// It must be called initially, before any other method is called.
	SetIndexFunc(f func(T) *int)
	// Changed restores the heap invariant after the item's value
	// has been changed. Call this method after modifying the value
	// of the item. If the item has been deleted or the heap has been
	// cleared, Changed does nothing.
	//
	// The Heap must have an index function.
	Changed(v T)
	// Delete removes the item with the given index from the heap.
	// If the item has already been deleted or the heap has been cleared,
	// Delete does nothing.
	//
	// The Heap must have an index function.
	Delete(v T)
}

// heapOrdered is a min-heap for ordered types.
type heapOrdered[T cmp.Ordered] struct {
	impl heapImpl[T]
}

// heapFunc is a min-heap for any type with a custom comparison function.
type heapFunc[T any] struct {
	impl    heapImpl[T]
	compare func(T, T) int
}

// heapImpl contains the data and provides shared implementation.
type heapImpl[T any] struct {
	values    []T
	indexes   []*int // from calls to indexFunc; updated in swap
	indexFunc func(T) *int
	mover     mover
}

// mover provides the up and down operations that differ between heap types.
type mover interface {
	up(i int)
	down(i int) bool
}

// newOrdered creates a new min-heap for ordered types.
func newOrdered[T cmp.Ordered]() Heap[T] {
	h := &heapOrdered[T]{}
	h.impl.mover = h
	return h
}

// NewFunc creates a new min-heap with a custom comparison function.
// The comparison function should return a negative value if a < b,
// zero if a == b, and a positive value if a > b.
func NewFunc[T any](compare func(T, T) int) Heap[T] {
	h := &heapFunc[T]{compare: compare}
	h.impl.mover = h
	return h
}

// SetIndexFunc sets a function that returns a pointer to an index field
// in the heap element.
func (h *heapOrdered[T]) SetIndexFunc(f func(T) *int) {
	h.impl.setIndexFunc(f)
}

// SetIndexFunc sets a function that returns a pointer to an index field
// in the heap element.
func (h *heapFunc[T]) SetIndexFunc(f func(T) *int) {
	h.impl.setIndexFunc(f)
}

func (h *heapImpl[T]) setIndexFunc(f func(T) *int) {
	h.indexFunc = f
	h.indexes = make([]*int, len(h.values))
}

// Insert adds an element to the heap.
func (h *heapOrdered[T]) Insert(value T) {
	h.impl.insert(value)
}

// Insert adds an element to the heap.
func (h *heapFunc[T]) Insert(value T) {
	h.impl.insert(value)
}

func (h *heapImpl[T]) insert(value T) {
	h.values = append(h.values, value)
	if h.indexes != nil {
		p := h.indexFunc(value)
		*p = len(h.indexes)
		h.indexes = append(h.indexes, p)
	}
	h.mover.up(len(h.values) - 1)
}

// InsertSlice adds all elements of s to the heap, then heapifies.
// The caller must not subsequently modify s.
func (h *heapOrdered[T]) InsertSlice(s []T) {
	h.impl.insertSlice(s)
}

// InsertSlice adds all elements of s to the heap, then heapifies.
// The caller must not subsequently modify s.
func (h *heapFunc[T]) InsertSlice(s []T) {
	h.impl.insertSlice(s)
}

func (h *heapImpl[T]) insertSlice(s []T) {
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
func (h *heapOrdered[T]) Min() T {
	return h.impl.min()
}

// Min returns the minimum element in the heap without removing it.
// It panics if the heap is empty.
func (h *heapFunc[T]) Min() T {
	return h.impl.min()
}

func (h *heapImpl[T]) min() T {
	if len(h.values) == 0 {
		panic("heap: Min called on empty heap")
	}
	return h.values[0]
}

// TakeMin removes and returns the minimum element from the heap.
// It panics if the heap is empty.
func (h *heapOrdered[T]) TakeMin() T {
	return h.impl.takeMin()
}

// TakeMin removes and returns the minimum element from the heap.
// It panics if the heap is empty.
func (h *heapFunc[T]) TakeMin() T {
	return h.impl.takeMin()
}

func (h *heapImpl[T]) takeMin() T {
	if len(h.values) == 0 {
		panic("heap: TakeMin called on empty heap")
	}
	min := h.values[0]
	h.deleteAt(0)
	return min
}

// ChangeMin replaces the minimum value in the heap with the given value.
// It panics if the heap is empty.
func (h *heapOrdered[T]) ChangeMin(v T) {
	h.impl.changeMin(v)
}

// ChangeMin replaces the minimum value in the heap with the given value.
// It panics if the heap is empty.
func (h *heapFunc[T]) ChangeMin(v T) {
	h.impl.changeMin(v)
}

func (h *heapImpl[T]) changeMin(v T) {
	if len(h.values) == 0 {
		panic("heap: ChangeMin called on empty heap")
	}
	h.values[0] = v
	h.mover.down(0)
}

func (h *heapImpl[T]) build() {
	n := len(h.values)
	for i := n/2 - 1; i >= 0; i-- {
		h.mover.down(i)
	}
}

// Clear removes all elements from the heap.
func (h *heapOrdered[T]) Clear() {
	h.impl.clear()
}

// Clear removes all elements from the heap.
func (h *heapFunc[T]) Clear() {
	h.impl.clear()
}

func (h *heapImpl[T]) clear() {
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
func (h *heapOrdered[T]) Len() int {
	return len(h.impl.values)
}

// Len returns the number of elements in the heap.
func (h *heapFunc[T]) Len() int {
	return len(h.impl.values)
}

// All returns an iterator over all elements in the heap
// in unspecified order.
func (h *heapOrdered[T]) All() iter.Seq[T] {
	return h.impl.all()
}

// All returns an iterator over all elements in the heap
// in unspecified order.
func (h *heapFunc[T]) All() iter.Seq[T] {
	return h.impl.all()
}

func (h *heapImpl[T]) all() iter.Seq[T] {
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
func (h *heapOrdered[T]) Drain() iter.Seq[T] {
	return h.impl.drain()
}

// Drain removes and returns the heap elements in sorted order,
// from smallest to largest.
//
// The result is undefined if the heap is changed during iteration.
func (h *heapFunc[T]) Drain() iter.Seq[T] {
	return h.impl.drain()
}

func (h *heapImpl[T]) drain() iter.Seq[T] {
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
func (h *heapOrdered[T]) Delete(v T) {
	h.impl.delete(v)
}

// Delete removes the item from the heap.
// If the item has already been deleted or the heap has been cleared,
// Delete does nothing.
//
// The Heap must have an index function.
func (h *heapFunc[T]) Delete(v T) {
	h.impl.delete(v)
}

func (h *heapImpl[T]) delete(v T) {
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

func (h *heapImpl[T]) deleteAt(i int) {
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
	if n != i && !h.mover.down(i) {
		h.mover.up(i)
	}
}

// Changed restores the heap invariant after the item's value has been changed.
// Call this method after modifying the value of the item.
// If the item has been deleted or the heap has been cleared, Changed does nothing.
//
// The Heap must have an index function.
func (h *heapOrdered[T]) Changed(v T) {
	h.impl.changed(v)
}

// Changed restores the heap invariant after the item's value has been changed.
// Call this method after modifying the value of the item.
// If the item has been deleted or the heap has been cleared, Changed does nothing.
//
// The Heap must have an index function.
func (h *heapFunc[T]) Changed(v T) {
	h.impl.changed(v)
}

func (h *heapImpl[T]) changed(v T) {
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
	if !h.mover.down(*ip) {
		h.mover.up(*ip)
	}
}

// up moves the element at index i up the heap until the heap invariant is restored.
func (h *heapOrdered[T]) up(i int) {
	for i > 0 {
		p := (i - 1) / 2 // parent
		if cmp.Compare(h.impl.values[i], h.impl.values[p]) >= 0 {
			break
		}
		h.impl.swap(p, i)
		i = p
	}
}

// up moves the element at index i up the heap until the heap invariant is restored.
func (h *heapFunc[T]) up(i int) {
	for i > 0 {
		p := (i - 1) / 2 // parent
		if h.compare(h.impl.values[i], h.impl.values[p]) >= 0 {
			break
		}
		h.impl.swap(p, i)
		i = p
	}
}

// down moves the element at index i down the heap until the heap invariant is restored.
// Returns true if the element moved.
func (h *heapOrdered[T]) down(i int) bool {
	n := len(h.impl.values)
	i0 := i
	for {
		lc := 2*i + 1
		if lc >= n {
			break
		}
		child := lc // left child
		if rc := lc + 1; rc < n && cmp.Compare(h.impl.values[rc], h.impl.values[lc]) < 0 {
			child = rc // right child is smaller
		}
		if cmp.Compare(h.impl.values[child], h.impl.values[i]) >= 0 {
			break
		}
		h.impl.swap(i, child)
		i = child
	}
	return i > i0
}

// down moves the element at index i down the heap until the heap invariant is restored.
// Returns true if the element moved.
func (h *heapFunc[T]) down(i int) bool {
	n := len(h.impl.values)
	i0 := i
	for {
		lc := 2*i + 1
		if lc >= n {
			break
		}
		child := lc // left child
		if rc := lc + 1; rc < n && h.compare(h.impl.values[rc], h.impl.values[lc]) < 0 {
			child = rc // right child is smaller
		}
		if h.compare(h.impl.values[child], h.impl.values[i]) >= 0 {
			break
		}
		h.impl.swap(i, child)
		i = child
	}
	return i > i0
}

func (h *heapImpl[T]) swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
	if h.indexes != nil {
		h.indexes[i], h.indexes[j] = h.indexes[j], h.indexes[i]
		*h.indexes[i] = i
		*h.indexes[j] = j
	}
}
