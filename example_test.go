package heap_test

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/jba/heap"
)

func ExampleHeap() {
	h := heap.New[int](cmp.Compare[int])

	// Insert elements.
	h.Init([]int{5, 3, 7, 1})

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

	h.Init([]int{5, 3, 7, 1})

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

	h := heap.NewIndexed(func(a, b *intWithIndex) int {
		return a.value - b.value
	}, func(v *intWithIndex, i int) { v.index = i })

	item1 := &intWithIndex{value: 5}
	item2 := &intWithIndex{value: 3}
	item3 := &intWithIndex{value: 7}
	item4 := &intWithIndex{value: 1}

	h.Init([]*intWithIndex{item1, item2, item3, item4})

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

	h := heap.NewIndexed(func(a, b *intWithIndex) int {
		return a.value - b.value
	}, func(v *intWithIndex, i int) { v.index = i })

	item1 := &intWithIndex{value: 5}
	item2 := &intWithIndex{value: 3}
	item3 := &intWithIndex{value: 7}

	h.Init([]*intWithIndex{item1, item2, item3})

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
	h.Init([]int{5, 3, 7, 1})

	// Iterate over all elements.
	sum := 0
	for v := range h.All() {
		sum += v
	}

	fmt.Printf("Total elements %d, sum %d\n", h.Len(), sum)

	// Output:
	// Total elements 4, sum 16
}

// ExampleHeap_topK demonstrates finding the K largest elements
// using a min-heap and ChangeMin.
// This is the "top K" algorithm.
func Example_topK() {
	// To find the K largest elements, use a min-heap of size K.
	// The heap's min is the smallest of the K largest seen so far.
	const K = 3
	data := []int{7, 2, 9, 1, 5, 8, 3, 6, 4, 10}

	h := heap.New(cmp.Compare[int])
	// Initialize the heap with the first K elements.
	// The heap will own data[:K], but we don't need to refer to that
	// part of the slice any more.
	h.Init(data[:K])

	// For remaining elements, replace the min if we find a larger value.
	for _, v := range data[K:] {
		if v > h.Min() {
			h.ChangeMin(v)
		}
	}

	// Drain the heap to get the K largest (in ascending order).
	fmt.Printf("%d largest elements:\n", K)
	for v := range h.Drain() {
		fmt.Println(v)
	}

	// Output:
	// 3 largest elements:
	// 8
	// 9
	// 10
}

func ExampleHeap_InsertAll() {
	h := heap.New(cmp.Compare[int])
	h.Init([]int{5, 3, 1})

	// Add more elements to the existing heap.
	h.InsertAll(slices.Values([]int{4, 2, 5, 6, 0}))

	for v := range h.Drain() {
		fmt.Println(v)
	}

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 5
	// 6
}

func Example_heapsort() {
	// To implement heapsort, first build a heap, then drain it.
	h := heap.New(cmp.Compare[int])
	h.Init([]int{7, 2, 9, 1, 5})
	for v := range h.Drain() {
		fmt.Println(v)
	}

	// Output:
	// 1
	// 2
	// 5
	// 7
	// 9
}
