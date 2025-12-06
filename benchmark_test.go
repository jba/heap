package heap

import (
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
				h.ExtractMin()
			}
		}
	})

	b.Run("InsertItem", func(b *testing.B) {
		for b.Loop() {
			h := New[int]()
			for _, n := range nums {
				h.InsertItem(n)
			}
			for h.Len() > 0 {
				h.ExtractMin()
			}
		}
	})
}
