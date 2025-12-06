package heap_test

import (
	"fmt"

	"github.com/jba/heap"
)

func ExampleHeap() {
	h := heap.New[int]()

	// Insert elements
	h.Insert(5)
	h.Insert(3)
	h.Insert(7)
	h.Insert(1)

	// Extract elements in sorted order
	fmt.Println(h.ExtractMin())
	fmt.Println(h.ExtractMin())
	fmt.Println(h.ExtractMin())
	fmt.Println(h.ExtractMin())

	// Output:
	// 1
	// 3
	// 5
	// 7
}

func ExampleHeap_Build() {
	h := heap.New[int]()

	// Insert many elements
	for _, v := range []int{5, 2, 8, 1, 9, 3, 7} {
		h.Insert(v)
	}

	// Build the heap explicitly to avoid cost on first Min/ExtractMin
	h.Build()

	fmt.Println(h.Min())

	// Output:
	// 1
}

func ExampleHeapFunc() {
	// Create a max-heap using a custom comparison function
	h := heap.NewFunc(func(a, b int) int {
		// Reverse the comparison for max-heap
		return b - a
	})

	h.Insert(5)
	h.Insert(3)
	h.Insert(7)
	h.Insert(1)

	// Extract maximum values
	fmt.Println(h.ExtractMin()) // "Min" extracts the element that compares smallest
	fmt.Println(h.ExtractMin())

	// Output:
	// 7
	// 5
}

func ExampleItem_Delete() {
	h := heap.New[int]()

	item1 := h.Insert(5)
	item2 := h.Insert(3)
	h.Insert(7)
	h.Insert(1)

	// Delete specific items
	item1.Delete()
	item2.Delete()

	// Remaining elements
	for h.Len() > 0 {
		fmt.Println(h.ExtractMin())
	}

	// Output:
	// 1
	// 7
}

func ExampleItem_Fix() {
	// In a real use case, you'd wrap your mutable data structure
	type mutableInt struct {
		value int
	}

	hm := heap.NewFunc(func(a, b *mutableInt) int {
		return a.value - b.value
	})

	val1 := &mutableInt{5}
	val2 := &mutableInt{3}
	val3 := &mutableInt{7}

	item1 := hm.Insert(val1)
	hm.Insert(val2)
	hm.Insert(val3)

	// Get the current min
	fmt.Println(hm.Min().value)

	// Modify val1's value (currently 5, make it smaller)
	val1.value = 1

	// Call Fix to restore heap invariant
	item1.Fix()

	// Now val1 should be the new minimum
	fmt.Println(hm.Min().value)

	// Output:
	// 3
	// 1
}

func ExampleHeap_All() {
	h := heap.New[int]()

	h.Insert(5)
	h.Insert(3)
	h.Insert(7)
	h.Insert(1)

	// Iterate over all elements
	// First element is guaranteed to be the minimum
	first := true
	for v := range h.All() {
		if first {
			fmt.Printf("Min: %d\n", v)
			first = false
		}
	}

	fmt.Printf("Total elements: %d\n", h.Len())

	// Output:
	// Min: 1
	// Total elements: 4
}
