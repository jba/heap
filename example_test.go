package heap_test

import (
	"cmp"
	"fmt"

	"github.com/jba/heap"
)

func ExampleNewFunc() {
	h := heap.NewFunc[int](cmp.Compare[int])

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
	h := heap.NewFunc(func(a, b int) int {
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

	h := heap.NewFunc(func(a, b *intWithIndex) int {
		return a.value - b.value
	})
	h.SetIndexFunc(func(v *intWithIndex) *int { return &v.index })

	item1 := &intWithIndex{value: 5}
	item2 := &intWithIndex{value: 3}
	item3 := &intWithIndex{value: 7}
	item4 := &intWithIndex{value: 1}

	h.InsertSlice([]*intWithIndex{item1, item2, item3, item4})

	// Delete specific items.
	h.Delete(item1)
	h.Delete(item2)

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

	h := heap.NewFunc(func(a, b *intWithIndex) int {
		return a.value - b.value
	})
	h.SetIndexFunc(func(v *intWithIndex) *int { return &v.index })

	item1 := &intWithIndex{value: 5}
	item2 := &intWithIndex{value: 3}
	item3 := &intWithIndex{value: 7}

	h.InsertSlice([]*intWithIndex{item1, item2, item3})

	// Get the current min.
	fmt.Println(h.Min().value)

	// Modify item1's value (currently 5, make it smaller).
	item1.value = 1

	// Call Changed to restore heap invariant.
	h.Changed(item1)

	// Now item1 should be the new minimum.
	fmt.Println(h.Min().value)

	// Output:
	// 3
	// 1
}

// ExampleHeap_ChangeMin demonstrates finding the K largest elements
// using a min-heap and ChangeMin.
// This is commonly known as the "top K" algorithm.
// func ExampleHeap_ChangeMin() {
// 	// To find the K largest elements, use a min-heap of size K.
// 	// The heap's min is the smallest of the K largest seen so far.
// 	h := heap.NewOrdered[int]()

// 	data := []int{7, 2, 9, 1, 5, 8, 3, 6, 4, 10}
// 	k := 3

// 	// Insert first K elements.
// 	h.InsertSlice(data[:k])

// 	// For remaining elements, replace the min if we find a larger value.
// 	for _, v := range data[k:] {
// 		if v > h.Min() {
// 			h.ChangeMin(v)
// 		}
// 	}

// 	// Drain the heap to get the K largest (in ascending order).
// 	fmt.Println("3 largest elements:")
// 	for v := range h.Drain() {
// 		fmt.Println(v)
// 	}

// 	// Output:
// 	// 3 largest elements:
// 	// 8
// 	// 9
// 	// 10
// }

// Example_kSmallest demonstrates finding the K smallest elements
// using a max-heap and ChangeMin.
func Example_kSmallest() {
	// To find the K smallest elements, use a max-heap of size K.
	// The heap's "min" (actually max) is the largest of the K smallest seen so far.
	h := heap.NewFunc(func(a, b int) int {
		return b - a // Reverse comparison for max-heap.
	})

	data := []int{7, 2, 9, 1, 5, 8, 3, 6, 4, 10}
	k := 3

	// Insert first K elements.
	h.InsertSlice(data[:k])

	// For remaining elements, replace the max if we find a smaller value.
	for _, v := range data[k:] {
		if v < h.Min() {
			h.ChangeMin(v)
		}
	}

	// Drain the heap to get the K smallest (in descending order).
	fmt.Println("3 smallest elements:")
	for v := range h.Drain() {
		fmt.Println(v)
	}

	// Output:
	// 3 smallest elements:
	// 3
	// 2
	// 1
}

// func ExampleHeap_All() {
// 	h := heap.NewOrdered[int]()
// 	h.InsertSlice([]int{5, 3, 7, 1})

// 	// Iterate over all elements.
// 	sum := 0
// 	for v := range h.All() {
// 		sum += v
// 	}

// 	fmt.Printf("Total elements %d, sum %d\n", h.Len(), sum)

// 	// Output:
// 	// Total elements 4, sum 16
// }
