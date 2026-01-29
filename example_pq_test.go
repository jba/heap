// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example demonstrates a priority queue built using the heap package.
package heap_test

import (
	"cmp"
	"fmt"

	"github.com/jba/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.
	index int // The index of the item in the heap.
}

// This example creates a priority queue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func Example_priorityQueue() {
	// Create a priority queue with highest priority first.
	// Since Heap is a min-heap, we reverse the comparison.
	pq := heap.NewIndexed(func(a, b *Item) int {
		return cmp.Compare(b.priority, a.priority)
	}, func(item *Item, i int) { item.index = i })

	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Add the items to the priority queue.
	for value, priority := range items {
		pq.Insert(&Item{
			value:    value,
			priority: priority,
		})
	}

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	pq.Insert(item)

	// Change the item's priority.
	item.priority = 5
	pq.Changed(item.index)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := pq.TakeMin()
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
	// Output:
	// 05:orange 04:pear 03:banana 02:apple
}
