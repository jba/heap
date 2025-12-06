// Package heap provides min-heap data structures.
package heap

import (
	"cmp"
	"iter"
)

// Heap is a min-heap for ordered types.
type Heap[T cmp.Ordered] struct {
	impl *heapImpl[T]
}

// HeapFunc is a min-heap for any type with a custom comparison function.
type HeapFunc[T any] struct {
	impl *heapImpl[T]
}

// Item represents an element in the heap and can be used to delete or fix it.
type Item struct {
	index *int
	heap  heapInterface
}

// heapInterface allows Item to call back into either Heap or HeapFunc.
type heapInterface interface {
	deleteItem(indexPtr *int)
	fixItem(indexPtr *int)
}

// entry stores a value and its index pointer.
type entry[T any] struct {
	value T
	index *int // shared with the Item; updated when the entry moves in the heap
}

// heapImpl is the internal implementation shared by Heap and HeapFunc.
type heapImpl[T any] struct {
	data    []entry[T]
	compare func(T, T) int
	built   bool // true if the heap invariant is currently maintained
}

// New creates a new min-heap for ordered types.
func New[T cmp.Ordered]() *Heap[T] {
	impl := &heapImpl[T]{
		data:    make([]entry[T], 0),
		compare: cmp.Compare[T],
	}
	return &Heap[T]{impl: impl}
}

// NewFunc creates a new min-heap with a custom comparison function.
// The comparison function should return a negative value if a < b,
// zero if a == b, and a positive value if a > b.
func NewFunc[T any](compare func(T, T) int) *HeapFunc[T] {
	impl := &heapImpl[T]{
		data:    make([]entry[T], 0),
		compare: compare,
	}
	return &HeapFunc[T]{impl: impl}
}

// Insert adds an element to the heap and returns an Item that can be used
// to delete or fix the element later.
//
// Before the first call to Min or ExtractMin, Insert simply appends to an
// internal slice without maintaining the heap invariant. Call Build explicitly
// if you want to ensure the heap is built after a batch of insertions.
func (h *Heap[T]) Insert(value T) Item {
	return h.impl.insert(value, h.impl)
}

// Insert adds an element to the heap and returns an Item that can be used
// to delete or fix the element later.
//
// Before the first call to Min or ExtractMin, Insert simply appends to an
// internal slice without maintaining the heap invariant. Call Build explicitly
// if you want to ensure the heap is built after a batch of insertions.
func (h *HeapFunc[T]) Insert(value T) Item {
	return h.impl.insert(value, h.impl)
}

func (h *heapImpl[T]) insert(value T, heap heapInterface) Item {
	idx := new(int)
	*idx = len(h.data)

	e := entry[T]{
		value: value,
		index: idx,
	}

	h.data = append(h.data, e)

	// If heap has already been built, maintain heap invariant
	if h.built {
		h.up(len(h.data) - 1)
	}

	return Item{
		index: idx,
		heap:  heap,
	}
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
	if !h.built {
		h.build()
	}
	if len(h.data) == 0 {
		panic("heap: Min called on empty heap")
	}
	return h.data[0].value
}

// ExtractMin removes and returns the minimum element from the heap.
// Panics if the heap is empty.
//
// The first call to ExtractMin builds the heap if it hasn't been built yet.
func (h *Heap[T]) ExtractMin() T {
	return h.impl.extractMin()
}

// ExtractMin removes and returns the minimum element from the heap.
// Panics if the heap is empty.
//
// The first call to ExtractMin builds the heap if it hasn't been built yet.
func (h *HeapFunc[T]) ExtractMin() T {
	return h.impl.extractMin()
}

func (h *heapImpl[T]) extractMin() T {
	if !h.built {
		h.build()
	}
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

func (h *heapImpl[T]) build() {
	n := len(h.data)
	// heapify: start from the last non-leaf node and sift down
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i)
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
	// Invalidate all outstanding Items
	for i := range h.data {
		*h.data[i].index = -1
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

// All returns an iterator over all elements in the heap.
// The first element yielded is guaranteed to be the minimum.
// The order of other elements is unspecified.
//
// If the heap hasn't been built yet, All builds it before iterating.
func (h *Heap[T]) All() iter.Seq[T] {
	return h.impl.all()
}

// All returns an iterator over all elements in the heap.
// The first element yielded is guaranteed to be the minimum.
// The order of other elements is unspecified.
//
// If the heap hasn't been built yet, All builds it before iterating.
func (h *HeapFunc[T]) All() iter.Seq[T] {
	return h.impl.all()
}

func (h *heapImpl[T]) all() iter.Seq[T] {
	return func(yield func(T) bool) {
		if !h.built && len(h.data) > 0 {
			h.build()
		}
		for _, e := range h.data {
			if !yield(e.value) {
				return
			}
		}
	}
}

// Delete removes this item from the heap.
// If the item has already been deleted or the heap has been cleared,
// Delete does nothing.
func (item Item) Delete() {
	if item.index == nil || *item.index < 0 {
		return // already deleted
	}
	item.heap.deleteItem(item.index)
}

// Fix restores the heap invariant after the item's value has been changed.
// Call this method after modifying the value of the element that this Item represents.
// If the item has been deleted or the heap has been cleared, Fix does nothing.
func (item Item) Fix() {
	if item.index == nil || *item.index < 0 {
		return // deleted item
	}
	item.heap.fixItem(item.index)
}

func (h *heapImpl[T]) deleteItem(indexPtr *int) {
	if !h.built {
		h.build()
	}
	// After building, read the updated index
	i := *indexPtr
	if i < 0 || i >= len(h.data) {
		return
	}
	h.deleteAt(i)
}

func (h *heapImpl[T]) fixItem(indexPtr *int) {
	if !h.built {
		h.build()
	}
	// After building, read the updated index
	i := *indexPtr
	if i < 0 || i >= len(h.data) {
		return
	}
	// Try moving down first, then up (only one will actually move the element)
	if !h.down(i) {
		h.up(i)
	}
}

// deleteAt removes the element at index i and restores the heap invariant.
func (h *heapImpl[T]) deleteAt(i int) {
	// Mark as deleted
	*h.data[i].index = -1

	n := len(h.data) - 1
	if n != i {
		h.swap(i, n)
		h.data = h.data[:n]
		// Restore heap invariant
		if !h.down(i) {
			h.up(i)
		}
	} else {
		h.data = h.data[:n]
	}
}

// up moves the element at index i up the heap until the heap invariant is restored.
func (h *heapImpl[T]) up(i int) {
	for {
		if i == 0 {
			break
		}
		p := (i - 1) / 2 // parent
		if h.compare(h.data[i].value, h.data[p].value) >= 0 {
			break
		}
		h.swap(p, i)
		i = p
	}
}

// down moves the element at index i down the heap until the heap invariant is restored.
// Returns true if the element moved.
func (h *heapImpl[T]) down(i int) bool {
	n := len(h.data)
	i0 := i
	for {
		// Find smallest child.
		lc := 2*i + 1
		if lc >= n {
			break
		}
		child := lc // left child
		if rc := lc + 1; rc < n && h.compare(h.data[rc].value, h.data[lc].value) < 0 {
			child = rc // right child is smaller
		}
		if h.compare(h.data[child].value, h.data[i].value) >= 0 {
			break
		}
		// Smaller child is less than parent.
		h.swap(i, child)
		i = child
	}
	return i > i0
}

// swap exchanges the elements at indices i and j and updates their index pointers.
func (h *heapImpl[T]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	*h.data[i].index = i
	*h.data[j].index = j
}
