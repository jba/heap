// Package heap provides min-heap data structures.
package heap

import (
	"cmp"
	"iter"
)

// Heap is a min-heap for ordered types.
type Heap[T cmp.Ordered] struct {
	impl heapImpl[T]
}

// HeapFunc is a min-heap for any type with a custom comparison function.
type HeapFunc[T any] struct {
	impl    heapImpl[T]
	compare func(T, T) int
}

// heapImpl contains the data and provides shared implementation.
// It uses the mover interface to call type-specific up/down methods.
type heapImpl[T any] struct {
	data  []entry[T]
	built bool // true if the heap invariant is currently maintained
	mover mover
}

// mover provides the up and down operations that differ between Heap and HeapFunc.
type mover interface {
	up(i int)
	down(i int) bool
}

// Handle represents an element in the heap and can be used to delete or modify it.
type Handle struct {
	index *int
	iface itemInterface
}

// itemInterface allows Item to call back into the heap implementation.
type itemInterface interface {
	deleteItem(index *int)
	changedItem(index *int)
}

// entry stores a value and its index pointer.
type entry[T any] struct {
	value T
	index *int // shared with the Item; updated when the entry moves in the heap
}

// New creates a new min-heap for ordered types.
func New[T cmp.Ordered]() *Heap[T] {
	h := &Heap[T]{}
	h.impl.mover = h
	return h
}

// NewFunc creates a new min-heap with a custom comparison function.
// The comparison function should return a negative value if a < b,
// zero if a == b, and a positive value if a > b.
func NewFunc[T any](compare func(T, T) int) *HeapFunc[T] {
	h := &HeapFunc[T]{compare: compare}
	h.impl.mover = h
	return h
}

// Insert adds an element to the heap.
//
// Before the first call to other methods such as Min or ExtractMin, Insert simply
// appends to an internal slice without maintaining the heap invariant. Call Build
// explicitly if you want to ensure the heap is built after a batch of insertions.
func (h *Heap[T]) Insert(value T) {
	h.impl.insert(entry[T]{value: value})
}

// Insert adds an element to the heap.
//
// Before the first call to other methods such as Min or ExtractMin, Insert simply
// appends to an internal slice without maintaining the heap invariant. Call Build
// explicitly if you want to ensure the heap is built after a batch of insertions.
func (h *HeapFunc[T]) Insert(value T) {
	h.impl.insert(entry[T]{value: value})
}

func (h *heapImpl[T]) insert(e entry[T]) {
	h.data = append(h.data, e)
	if h.built {
		h.mover.up(len(h.data) - 1)
	}
}

// InsertHandle adds an element to the heap and returns an Item that can be used
// to delete or adjust the element later.
//
// Before the first call to other methods such as Min or ExtractMin, InsertHandle simply
// appends to an internal slice without maintaining the heap invariant. Call Build
// explicitly if you want to ensure the heap is built after a batch of insertions.
func (h *Heap[T]) InsertHandle(value T) Handle {
	return h.impl.insertHandle(entry[T]{value: value})
}

// InsertHandle adds an element to the heap and returns an Item that can be used
// to delete or adjust the element later.
//
// Before the first call to other methods such as Min or ExtractMin, InsertHandle simply
// appends to an internal slice without maintaining the heap invariant. Call Build
// explicitly if you want to ensure the heap is built after a batch of insertions.
func (h *HeapFunc[T]) InsertHandle(value T) Handle {
	return h.impl.insertHandle(entry[T]{value: value})
}

func (h *heapImpl[T]) insertHandle(e entry[T]) Handle {
	idx := new(int)
	*idx = len(h.data)
	e.index = idx
	h.data = append(h.data, e)
	if h.built {
		h.mover.up(len(h.data) - 1)
	}
	return Handle{index: idx, iface: h}
}

// Min returns the minimum element in the heap without removing it.
// Panics if the heap is empty.
//
// The first call to Min builds the heap if it hasn't been built yet.
func (h *Heap[T]) Min() T {
	return h.impl.min()
}

// Min returns the minimum element in the heap without removing it.
// Panics if the heap is empty.
//
// The first call to Min builds the heap if it hasn't been built yet.
func (h *HeapFunc[T]) Min() T {
	return h.impl.min()
}

func (h *heapImpl[T]) min() T {
	h.ensureBuilt()
	if len(h.data) == 0 {
		panic("heap: Min called on empty heap")
	}
	return h.data[0].value
}

// TakeMin removes and returns the minimum element from the heap.
// Panics if the heap is empty.
//
// The first call to TakeMin builds the heap if it hasn't been built yet.
func (h *Heap[T]) TakeMin() T {
	return h.impl.takeMin()
}

// TakeMin removes and returns the minimum element from the heap.
// Panics if the heap is empty.
//
// The first call to TakeMin builds the heap if it hasn't been built yet.
func (h *HeapFunc[T]) TakeMin() T {
	return h.impl.takeMin()
}

func (h *heapImpl[T]) takeMin() T {
	h.ensureBuilt()
	if len(h.data) == 0 {
		panic("heap: ExtractMin called on empty heap")
	}
	min := h.data[0].value
	h.deleteAt(0)
	return min
}

// Build rebuilds the heap in O(n) time.
// Call this after inserting multiple elements to avoid the cost of building
// the heap on the first call to Min or ExtractMin.
func (h *Heap[T]) Build() {
	h.impl.build()
}

// Build rebuilds the heap in O(n) time.
// Call this after inserting multiple elements to avoid the cost of building
// the heap on the first call to Min or ExtractMin.
func (h *HeapFunc[T]) Build() {
	h.impl.build()
}

func (h *heapImpl[T]) ensureBuilt() {
	if !h.built {
		h.build()
	}
}

func (h *heapImpl[T]) build() {
	n := len(h.data)
	for i := n/2 - 1; i >= 0; i-- {
		h.mover.down(i)
	}
	h.built = true
}

// Clear removes all elements from the heap.
func (h *Heap[T]) Clear() {
	h.impl.clear()
}

// Clear removes all elements from the heap.
func (h *HeapFunc[T]) Clear() {
	h.impl.clear()
}

func (h *heapImpl[T]) clear() {
	for i := range h.data {
		if h.data[i].index != nil {
			*h.data[i].index = -1
		}
	}
	h.data = h.data[:0]
	h.built = false
}

// Len returns the number of elements in the heap.
func (h *Heap[T]) Len() int {
	return len(h.impl.data)
}

// Len returns the number of elements in the heap.
func (h *HeapFunc[T]) Len() int {
	return len(h.impl.data)
}

// All returns an iterator over all elements in the heap
// in unspecified order.
func (h *Heap[T]) All() iter.Seq[T] {
	return h.impl.all()
}

// All returns an iterator over all elements in the heap
// in unspecified order.
func (h *HeapFunc[T]) All() iter.Seq[T] {
	return h.impl.all()
}

func (h *heapImpl[T]) all() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, e := range h.data {
			if !yield(e.value) {
				return
			}
		}
	}
}

// Drain removes and returns the heap elements in sorted order,
// from smallest to largest.
func (h *Heap[T]) Drain() iter.Seq[T] {
	return h.impl.drain()
}

// Drain removes and returns the heap elements in sorted order,
// from smallest to largest.
func (h *HeapFunc[T]) Drain() iter.Seq[T] {
	return h.impl.drain()
}

func (h *heapImpl[T]) drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for len(h.data) > 0 {
			if !yield(h.takeMin()) {
				return
			}
		}
	}
}

// Delete removes this item from the heap.
// If the item has already been deleted or the heap has been cleared,
// Delete does nothing.
func (item Handle) Delete() {
	if item.index == nil || *item.index < 0 {
		return // already deleted
	}
	item.iface.deleteItem(item.index)
}

// Changed restores the heap invariant after the item's value has been changed.
// Call this method after modifying the value of the element that this Item represents.
// If the item has been deleted or the heap has been cleared, Changed does nothing.
func (item Handle) Changed() {
	if item.index == nil || *item.index < 0 {
		return // deleted item
	}
	item.iface.changedItem(item.index)
}

func (h *heapImpl[T]) deleteItem(index *int) {
	h.ensureBuilt()
	i := *index
	if i < 0 || i >= len(h.data) {
		return
	}
	h.deleteAt(i)
}

func (h *heapImpl[T]) changedItem(index *int) {
	h.ensureBuilt()
	i := *index
	if i < 0 || i >= len(h.data) {
		return
	}
	if !h.mover.down(i) {
		h.mover.up(i)
	}
}

func (h *heapImpl[T]) deleteAt(i int) {
	if h.data[i].index != nil {
		*h.data[i].index = -1
	}
	n := len(h.data) - 1
	if n != i {
		h.swap(i, n)
		h.data = h.data[:n]
		if !h.mover.down(i) {
			h.mover.up(i)
		}
	} else {
		h.data = h.data[:n]
	}
}

// up moves the element at index i up the heap until the heap invariant is restored.
func (h *Heap[T]) up(i int) {
	for i > 0 {
		p := (i - 1) / 2 // parent
		if cmp.Compare(h.impl.data[i].value, h.impl.data[p].value) >= 0 {
			break
		}
		h.impl.swap(p, i)
		i = p
	}
}

// up moves the element at index i up the heap until the heap invariant is restored.
func (h *HeapFunc[T]) up(i int) {
	for i > 0 {
		p := (i - 1) / 2 // parent
		if h.compare(h.impl.data[i].value, h.impl.data[p].value) >= 0 {
			break
		}
		h.impl.swap(p, i)
		i = p
	}
}

// down moves the element at index i down the heap until the heap invariant is restored.
// Returns true if the element moved.
func (h *Heap[T]) down(i int) bool {
	data := h.impl.data
	n := len(data)
	i0 := i
	for {
		lc := 2*i + 1
		if lc >= n {
			break
		}
		child := lc // left child
		if rc := lc + 1; rc < n && cmp.Compare(data[rc].value, data[lc].value) < 0 {
			child = rc // right child is smaller
		}
		if cmp.Compare(data[child].value, data[i].value) >= 0 {
			break
		}
		h.impl.swap(i, child)
		i = child
	}
	return i > i0
}

// down moves the element at index i down the heap until the heap invariant is restored.
// Returns true if the element moved.
func (h *HeapFunc[T]) down(i int) bool {
	data := h.impl.data
	n := len(data)
	i0 := i
	for {
		lc := 2*i + 1
		if lc >= n {
			break
		}
		child := lc // left child
		if rc := lc + 1; rc < n && h.compare(data[rc].value, data[lc].value) < 0 {
			child = rc // right child is smaller
		}
		if h.compare(data[child].value, data[i].value) >= 0 {
			break
		}
		h.impl.swap(i, child)
		i = child
	}
	return i > i0
}

func (h *heapImpl[T]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	if h.data[i].index != nil {
		*h.data[i].index = i
	}
	if h.data[j].index != nil {
		*h.data[j].index = j
	}
}
