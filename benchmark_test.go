package heap

import (
	"cmp"
	"math/rand"
	"slices"
	"testing"
)

func BenchmarkHeapsort(b *testing.B) {
	b.Run("int", func(b *testing.B) {
		nums := make([]int, 1000)
		for i := range nums {
			nums[i] = rand.Int()
		}
		b.ResetTimer()
		for b.Loop() {
			h := New(cmp.Compare[int])
			h.InsertSlice(slices.Clone(nums))
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})

	b.Run("struct", func(b *testing.B) {
		nums := make([]*intIndexed, 1000)
		for i := range nums {
			nums[i] = &intIndexed{value: rand.Int()}
		}
		b.ResetTimer()
		for b.Loop() {
			h := New(func(a, b *intIndexed) int { return cmp.Compare(a.value, b.value) })
			h.InsertSlice(slices.Clone(nums))
			for h.Len() > 0 {
				h.TakeMin()
			}
		}
	})
}
