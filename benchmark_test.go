package heap

import (
	"cmp"
	"math/rand/v2"
	"slices"
	"testing"
)

type benchTask struct {
	priority int
	index    int
}

func BenchmarkHeapsort(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = rand.Int()
	}
	b.Run("Int", func(b *testing.B) {
		for b.Loop() {
			h := New(cmp.Compare[int])
			h.Init(slices.Clone(nums))
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})
	b.Run("Ordered", func(b *testing.B) {
		for b.Loop() {
			h := newOrderedHeap[int]()
			h.Init(slices.Clone(nums))
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})

	b.Run("Struct", func(b *testing.B) {
		nums := make([]*intIndexed, 1000)
		for i := range nums {
			nums[i] = &intIndexed{value: rand.Int()}
		}
		b.ResetTimer()
		for b.Loop() {
			h := New(func(a, b *intIndexed) int { return cmp.Compare(a.value, b.value) })
			h.Init(slices.Clone(nums))
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})
}

func BenchmarkPriorityQueue(b *testing.B) {
	cmpTask := func(a, b *benchTask) int { return cmp.Compare(a.priority, b.priority) }

	// Pre-generate all random numbers for deterministic iterations
	const nTasks = 100
	const nRounds = 50

	initialPriorities := make([]int, nTasks)
	for i := range initialPriorities {
		initialPriorities[i] = rand.IntN(1000)
	}

	// For each round: 3 task indices and 3 new priorities for Changed
	changeTaskIdx := make([]int, nRounds*3)
	changePriority := make([]int, nRounds*3)
	for i := range changeTaskIdx {
		changeTaskIdx[i] = rand.IntN(nTasks)
		changePriority[i] = rand.IntN(1000)
	}

	for b.Loop() {
		h := NewIndexed(cmpTask, func(t *benchTask, i int) { t.index = i })

		// Pool of tasks we can add/remove/modify
		tasks := make([]*benchTask, nTasks)
		for i := range tasks {
			tasks[i] = &benchTask{priority: initialPriorities[i]}
		}

		// Simulate priority queue workload
		for round := range nRounds {
			// Add some tasks
			for i := range 10 {
				h.Insert(tasks[(round*10+i)%len(tasks)])
			}

			// Process highest priority tasks
			for range 5 {
				if h.Len() > 0 {
					h.TakeMin()
				}
			}

			// Change priority of some tasks still in heap
			for j := range 3 {
				idx := changeTaskIdx[round*3+j]
				t := tasks[idx]
				if t.index >= 0 && t.index < h.Len() {
					t.priority = changePriority[round*3+j]
					h.Changed(t.index)
				}
			}

			// Delete a random task from heap
			if h.Len() > 1 {
				h.Delete(1)
			}
		}

		// Drain remaining
		for h.Len() > 0 {
			h.TakeMin()
		}
	}
}

func BenchmarkTopK(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = rand.Int()
	}

	const k = 5

	b.Run("kind=Ordered", func(b *testing.B) {
		for b.Loop() {
			h := newOrderedHeap[int]()

			// Insert first k elements
			h.Init(data[:k])

			// For remaining elements, replace min if we find a larger value
			for _, v := range data[k:] {
				if v > h.Min() {
					h.ChangeMin(v)
				}
			}
		}
	})
	b.Run("kind=int", func(b *testing.B) {
		for b.Loop() {
			h := newIntHeap()

			// Insert first k elements
			h.Init(data[:k])

			// For remaining elements, replace min if we find a larger value
			for _, v := range data[k:] {
				if v > h.Min() {
					h.ChangeMin(v)
				}
			}
		}
	})
	b.Run("kind=Grafana", func(b *testing.B) {
		for b.Loop() {
			// Copy first k elements and heapify
			h := make([]int, k)
			copy(h, data[:k])
			grafanaHeapify(h)

			// For remaining elements, replace min if we find a larger value
			for _, v := range data[k:] {
				if v > grafanaMin(h) {
					grafanaChangeMin(h, v)
				}
			}
		}
	})
	b.Run("kind=Heap", func(b *testing.B) {
		for b.Loop() {
			h := New[int](cmp.Compare[int])

			// Insert first k elements
			h.Init(data[:k])

			// For remaining elements, replace min if we find a larger value
			for _, v := range data[k:] {
				if v > h.Min() {
					h.ChangeMin(v)
				}
			}
		}
	})
}
