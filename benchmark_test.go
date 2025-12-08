package heap

import (
	"cmp"
	"math/rand"
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = rand.Int()
	}

	b.Run("Insert", func(b *testing.B) {
		for b.Loop() {
			h := New[int]()
			for _, n := range nums {
				h.Insert(n)
			}
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})

	b.Run("InsertItem", func(b *testing.B) {
		for b.Loop() {
			h := New[int]()
			for _, n := range nums {
				h.InsertHandle(n)
			}
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})
}

func BenchmarkHeapVsHeapFunc(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = rand.Int()
	}

	b.Run("Heap", func(b *testing.B) {
		for b.Loop() {
			h := New[int]()
			for _, n := range nums {
				h.Insert(n)
			}
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})

	b.Run("HeapFunc", func(b *testing.B) {
		for b.Loop() {
			h := NewFunc(cmp.Compare[int])
			for _, n := range nums {
				h.Insert(n)
			}
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})
}
