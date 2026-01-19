package heap

import (
	"cmp"
	"slices"
	"testing"
)

func TestHeapBasicOperations(t *testing.T) {
	h := New(cmp.Compare[int])

	// Test empty heap
	if h.Len() != 0 {
		t.Errorf("new heap should have length 0, got %d", h.Len())
	}

	// Test Insert and Len
	h.Insert(5)
	h.Insert(3)
	h.Insert(7)
	h.Insert(1)

	if h.Len() != 4 {
		t.Errorf("heap should have length 4, got %d", h.Len())
	}

	// Test Min
	if min := h.Min(); min != 1 {
		t.Errorf("Min() = %d, want 1", min)
	}

	// Min should not remove element
	if h.Len() != 4 {
		t.Errorf("Min() should not remove element, len = %d", h.Len())
	}

	// Test ExtractMin
	if min := h.TakeMin(); min != 1 {
		t.Errorf("ExtractMin() = %d, want 1", min)
	}
	if h.Len() != 3 {
		t.Errorf("after ExtractMin, len should be 3, got %d", h.Len())
	}

	if min := h.TakeMin(); min != 3 {
		t.Errorf("ExtractMin() = %d, want 3", min)
	}
	if min := h.TakeMin(); min != 5 {
		t.Errorf("ExtractMin() = %d, want 5", min)
	}
	if min := h.TakeMin(); min != 7 {
		t.Errorf("ExtractMin() = %d, want 7", min)
	}

	if h.Len() != 0 {
		t.Errorf("heap should be empty, len = %d", h.Len())
	}
}

func TestHeapBuild(t *testing.T) {
	h := New(cmp.Compare[int])
	h.InsertSlice([]int{5, 2, 8, 1, 9, 3, 7})

	// Extract all elements - should come out in sorted order
	extracted := slices.Collect(h.Drain())

	expected := []int{1, 2, 3, 5, 7, 8, 9}
	if !slices.Equal(extracted, expected) {
		t.Errorf("extracted = %v, want %v", extracted, expected)
	}
}

func TestHeapFunc(t *testing.T) {
	// Create a max-heap by reversing the comparison
	h := New(func(a, b int) int {
		if a > b {
			return -1
		} else if a < b {
			return 1
		}
		return 0
	})
	h.InsertSlice([]int{5, 3, 7, 1})

	// Should extract in descending order
	if max := h.TakeMin(); max != 7 {
		t.Errorf("ExtractMin() = %d, want 7", max)
	}
	if max := h.TakeMin(); max != 5 {
		t.Errorf("ExtractMin() = %d, want 5", max)
	}
	if max := h.TakeMin(); max != 3 {
		t.Errorf("ExtractMin() = %d, want 3", max)
	}
	if max := h.TakeMin(); max != 1 {
		t.Errorf("ExtractMin() = %d, want 1", max)
	}
}

type intIndexed struct {
	value int
	index int
}

func TestItemDelete(t *testing.T) {
	h := New(func(a, b *intIndexed) int {
		return cmp.Compare(a.value, b.value)
	})
	h.SetIndexFunc(func(v *intIndexed, i int) { v.index = i })
	items := []*intIndexed{
		{value: 1},
		{value: 3},
		{value: 7},
		{value: 5},
	}
	h.InsertSlice(slices.Clone(items))

	if h.Len() != 4 {
		t.Fatalf("heap should have 4 elements, got %d", h.Len())
	}

	// Delete the middle element (value 3)
	h.Delete(items[1].index)
	if h.Len() != 3 {
		t.Errorf("after Delete, heap should have 3 elements, got %d", h.Len())
	}

	// Extract all remaining elements
	var got []int
	for e := range h.Drain() {
		got = append(got, e.value)
	}

	want := []int{1, 5, 7}
	if !slices.Equal(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}

	// Deleting already-deleted items (index -1) should be safe.
	for _, i := range items {
		h.Delete(i.index)
	}

	if g := h.Len(); g != 0 {
		t.Errorf("want zero len, got %d", g)
	}
}

func TestItemAdjust(t *testing.T) {
	h := New(func(a, b *intIndexed) int {
		return cmp.Compare(a.value, b.value)
	})
	h.SetIndexFunc(func(v *intIndexed, i int) { v.index = i })

	items := []*intIndexed{
		{value: 5},
		{value: 3},
		{value: 7},
		{value: 1},
		{value: 9},
	}
	// Save reference before InsertSlice reorders the slice.
	itemToModify := items[3]
	h.InsertSlice(items)

	// Modify the item (originally value 1) to 8.
	itemToModify.value = 8
	h.Changed(itemToModify.index)

	// Extract all elements - should still be in sorted order.
	var got []int
	for v := range h.Drain() {
		got = append(got, v.value)
	}

	want := []int{3, 5, 7, 8, 9}
	if !slices.Equal(got, want) {
		t.Errorf("after Adjust, got = %v, want %v", got, want)
	}
}

func TestClear(t *testing.T) {
	h := New(func(a, b *intIndexed) int {
		return cmp.Compare(a.value, b.value)
	})
	h.SetIndexFunc(func(v *intIndexed, i int) { v.index = i })

	items := []*intIndexed{
		{value: 5},
		{value: 3},
		{value: 7},
	}
	// Use Insert individually since we need to keep items valid after Clear.
	for _, item := range items {
		h.Insert(item)
	}

	h.Clear()

	if h.Len() != 0 {
		t.Errorf("after Clear, len should be 0, got %d", h.Len())
	}

	// Verify indices are set to -1 after clear.
	for _, item := range items {
		if item.index != -1 {
			t.Errorf("after Clear, item.index = %d, want -1", item.index)
		}
	}

	// Operations on items from cleared heap should be safe (index is -1).
	h.Delete(items[0].index)
	h.Changed(items[1].index)
}

func TestAll(t *testing.T) {
	h := New(cmp.Compare[int])
	h.InsertSlice([]int{5, 2, 8, 1, 9})

	// Collect all elements
	var collected []int
	for v := range h.All() {
		collected = append(collected, v)
	}

	if len(collected) != 5 {
		t.Errorf("All() yielded %d elements, want 5", len(collected))
	}

	// All original values should be present
	slices.Sort(collected)
	expected := []int{1, 2, 5, 8, 9}
	if !slices.Equal(collected, expected) {
		t.Errorf("All() values = %v, want %v", collected, expected)
	}
}

func TestAllEarlyBreak(t *testing.T) {
	h := New(cmp.Compare[int])

	for i := range 10 {
		h.Insert(i)
	}

	// Test that breaking early works
	count := 0
	for range h.All() {
		count++
		if count >= 3 {
			break
		}
	}

	if count != 3 {
		t.Errorf("broke after %d iterations, want 3", count)
	}

	// Heap should still be intact
	if h.Len() != 10 {
		t.Errorf("heap len = %d, want 10", h.Len())
	}
}

func TestPanicOnEmptyHeap(t *testing.T) {
	h := New(cmp.Compare[int])

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Min() on empty heap should panic")
		}
	}()

	h.Min()
}

func TestPanicOnEmptyExtractMin(t *testing.T) {
	h := New(cmp.Compare[int])

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ExtractMin() on empty heap should panic")
		}
	}()

	h.TakeMin()
}

func TestHeapWithStrings(t *testing.T) {
	h := New(cmp.Compare[string])
	h.InsertSlice([]string{"dog", "cat", "bird", "ant"})

	if min := h.Min(); min != "ant" {
		t.Errorf("Min() = %q, want %q", min, "ant")
	}

	extracted := slices.Collect(h.Drain())

	expected := []string{"ant", "bird", "cat", "dog"}
	if !slices.Equal(extracted, expected) {
		t.Errorf("extracted = %v, want %v", extracted, expected)
	}
}

func TestLargeHeap(t *testing.T) {
	h := New(cmp.Compare[int])

	// Insert 1000 elements in reverse order
	for i := 1000; i > 0; i-- {
		h.Insert(i)
	}

	// Extract all and verify they come out sorted
	prev := 0
	for curr := range h.Drain() {
		if curr <= prev {
			t.Errorf("extracted %d after %d, not in sorted order", curr, prev)
		}
		prev = curr
	}

	if prev != 1000 {
		t.Errorf("last extracted value = %d, want 1000", prev)
	}
}

func TestIndexFuncNoAllocs(t *testing.T) {
	h := New(func(a, b *intIndexed) int {
		return cmp.Compare(a.value, b.value)
	})
	h.SetIndexFunc(func(v *intIndexed, i int) { v.index = i })

	items := make([]*intIndexed, 100)
	for i := range items {
		items[i] = &intIndexed{value: i}
	}
	h.InsertSlice(items)

	// Verify Delete requires no allocation.
	allocs := testing.AllocsPerRun(5, func() {
		h.Delete(items[50].index)
		// Re-insert so we can delete again.
		items[50].value = 50
		h.Insert(items[50])
	})
	if allocs != 0 {
		t.Errorf("Delete: got %v allocs, want 0", allocs)
	}

	// Verify Changed requires no allocation.
	allocs = testing.AllocsPerRun(5, func() {
		items[25].value = 1000
		h.Changed(items[25].index)
		items[25].value = 25
		h.Changed(items[25].index)
	})
	if allocs != 0 {
		t.Errorf("Changed: got %v allocs, want 0", allocs)
	}
}

func TestInsertSlice(t *testing.T) {
	t.Run("Heap", func(t *testing.T) {
		h := New(cmp.Compare[int])
		data := []int{7, 2, 9, 1, 5}
		h.InsertSlice(data)

		got := slices.Collect(h.Drain())
		want := []int{1, 2, 5, 7, 9}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Heap takes ownership", func(t *testing.T) {
		h := New(cmp.Compare[int])
		data := []int{7, 2, 9, 1, 5}
		h.InsertSlice(data)

		// Modifying original slice should affect heap (ownership taken).
		data[0] = 100
		found := false
		for v := range h.All() {
			if v == 100 {
				found = true
				break
			}
		}
		if !found {
			t.Error("heap should have taken ownership of slice")
		}
	})

	t.Run("Heap appends to existing", func(t *testing.T) {
		h := New(cmp.Compare[int])
		h.Insert(10)
		h.Insert(20)

		h.InsertSlice([]int{5, 15, 25})

		got := slices.Collect(h.Drain())
		want := []int{5, 10, 15, 20, 25}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("HeapFunc without indexFunc", func(t *testing.T) {
		h := New(func(a, b int) int { return cmp.Compare(a, b) })
		data := []int{7, 2, 9, 1, 5}
		h.InsertSlice(data)

		got := slices.Collect(h.Drain())
		want := []int{1, 2, 5, 7, 9}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("HeapFunc with indexFunc", func(t *testing.T) {
		h := New(func(a, b *intIndexed) int {
			return cmp.Compare(a.value, b.value)
		})
		h.SetIndexFunc(func(v *intIndexed, i int) { v.index = i })

		items := []*intIndexed{
			{value: 7},
			{value: 2},
			{value: 9},
			{value: 1},
			{value: 5},
		}
		h.InsertSlice(items)

		// Verify indexes are set correctly.
		for _, item := range items {
			if item.index < 0 || item.index >= len(items) {
				t.Errorf("item with value %d has invalid index %d", item.value, item.index)
			}
		}

		// Delete an item to verify indexes work.
		h.Delete(items[1].index) // value 2

		var got []int
		for v := range h.Drain() {
			got = append(got, v.value)
		}
		want := []int{1, 5, 7, 9}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("HeapFunc with indexFunc appends to existing", func(t *testing.T) {
		h := New(func(a, b *intIndexed) int {
			return cmp.Compare(a.value, b.value)
		})
		h.SetIndexFunc(func(v *intIndexed, i int) { v.index = i })

		// Insert initial items.
		initial := []*intIndexed{{value: 10}, {value: 20}}
		for _, item := range initial {
			h.Insert(item)
		}

		// InsertSlice more items.
		more := []*intIndexed{{value: 5}, {value: 15}, {value: 25}}
		h.InsertSlice(more)

		// Verify all indexes are valid.
		all := append(initial, more...)
		for _, item := range all {
			if item.index < 0 || item.index >= len(all) {
				t.Errorf("item with value %d has invalid index %d", item.value, item.index)
			}
		}

		var got []int
		for v := range h.Drain() {
			got = append(got, v.value)
		}
		want := []int{5, 10, 15, 20, 25}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
