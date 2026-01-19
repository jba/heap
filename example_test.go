package heap_test

import (
	"cmp"
	"fmt"

	"github.com/jba/heap"
)

func ExampleHeap() {
	h := heap.New[int](cmp.Compare[int])

	// Insert elements.
	h.InsertSlice([]int{5, 3, 7, 1})

	// Extract elements in sorted order.
	fmt.Println(h.TakeMin())
	fmt.Println(h.TakeMin())
	fmt.Println(h.TakeMin())
	fmt.Println(h.TakeMin())

	// Output:
	// 1
	// 3
	// 5
	// 7
}

func Example_maxHeap() {
	// Create a max-heap using a custom comparison function.
	h := heap.New(func(a, b int) int {
		// Reverse the comparison for max-heap.
		return b - a
	})

	h.InsertSlice([]int{5, 3, 7, 1})

	// Extract maximum values.
	fmt.Println(h.TakeMin()) // "Min" extracts the element that compares smallest
	fmt.Println(h.TakeMin())

	// Output:
	// 7
	// 5
}

func Example_delete() {
	type intWithIndex struct {
		value int
		index int
	}

	h := heap.New(func(a, b *intWithIndex) int {
		return a.value - b.value
	})
	h.SetIndexFunc(func(v *intWithIndex, i int) { v.index = i })

	item1 := &intWithIndex{value: 5}
	item2 := &intWithIndex{value: 3}
	item3 := &intWithIndex{value: 7}
	item4 := &intWithIndex{value: 1}

	h.InsertSlice([]*intWithIndex{item1, item2, item3, item4})

	// Delete specific items by their index.
	h.Delete(item1.index)
	h.Delete(item2.index)

	// Remaining elements.
	for v := range h.Drain() {
		fmt.Println(v.value)
	}

	// Output:
	// 1
	// 7
}

func Example_changed() {
	type intWithIndex struct {
		value int
		index int
	}

	h := heap.New(func(a, b *intWithIndex) int {
		return a.value - b.value
	})
	h.SetIndexFunc(func(v *intWithIndex, i int) { v.index = i })

	item1 := &intWithIndex{value: 5}
	item2 := &intWithIndex{value: 3}
	item3 := &intWithIndex{value: 7}

	h.InsertSlice([]*intWithIndex{item1, item2, item3})

	// Get the current min.
	fmt.Println(h.Min().value)

	// Modify item1's value (currently 5, make it smaller).
	item1.value = 1

	// Call Changed to restore heap invariant.
	h.Changed(item1.index)

	// Now item1 should be the new minimum.
	fmt.Println(h.Min().value)

	// Output:
	// 3
	// 1
}

func ExampleHeap_All() {
	h := heap.New(cmp.Compare[int])
	h.InsertSlice([]int{5, 3, 7, 1})

	// Iterate over all elements.
	sum := 0
	for v := range h.All() {
		sum += v
	}

	fmt.Printf("Total elements %d, sum %d\n", h.Len(), sum)

	// Output:
	// Total elements 4, sum 16
}
