package heap

import (
	"slices"
	"testing"
)

func TestHeapBasicOperations(t *testing.T) {
	h := New[int]()

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
	if min := h.ExtractMin(); min != 1 {
		t.Errorf("ExtractMin() = %d, want 1", min)
	}
	if h.Len() != 3 {
		t.Errorf("after ExtractMin, len should be 3, got %d", h.Len())
	}

	if min := h.ExtractMin(); min != 3 {
		t.Errorf("ExtractMin() = %d, want 3", min)
	}
	if min := h.ExtractMin(); min != 5 {
		t.Errorf("ExtractMin() = %d, want 5", min)
	}
	if min := h.ExtractMin(); min != 7 {
		t.Errorf("ExtractMin() = %d, want 7", min)
	}

	if h.Len() != 0 {
		t.Errorf("heap should be empty, len = %d", h.Len())
	}
}

func TestHeapBuild(t *testing.T) {
	h := New[int]()

	// Insert several elements
	values := []int{5, 2, 8, 1, 9, 3, 7}
	for _, v := range values {
		h.Insert(v)
	}

	// Explicitly build the heap
	h.Build()

	// Extract all elements - should come out in sorted order
	var extracted []int
	for h.Len() > 0 {
		extracted = append(extracted, h.ExtractMin())
	}

	expected := []int{1, 2, 3, 5, 7, 8, 9}
	if !slices.Equal(extracted, expected) {
		t.Errorf("extracted = %v, want %v", extracted, expected)
	}
}

func TestHeapFunc(t *testing.T) {
	// Create a max-heap by reversing the comparison
	h := NewFunc(func(a, b int) int {
		if a > b {
			return -1
		} else if a < b {
			return 1
		}
		return 0
	})

	h.Insert(5)
	h.Insert(3)
	h.Insert(7)
	h.Insert(1)

	// Should extract in descending order
	if max := h.ExtractMin(); max != 7 {
		t.Errorf("ExtractMin() = %d, want 7", max)
	}
	if max := h.ExtractMin(); max != 5 {
		t.Errorf("ExtractMin() = %d, want 5", max)
	}
	if max := h.ExtractMin(); max != 3 {
		t.Errorf("ExtractMin() = %d, want 3", max)
	}
	if max := h.ExtractMin(); max != 1 {
		t.Errorf("ExtractMin() = %d, want 1", max)
	}
}

func TestItemDelete(t *testing.T) {
	h := New[int]()

	item1 := h.Insert(5)
	item2 := h.Insert(3)
	item3 := h.Insert(7)
	item4 := h.Insert(1)

	if h.Len() != 4 {
		t.Fatalf("heap should have 4 elements, got %d", h.Len())
	}

	// Delete the middle element
	item2.Delete()
	if h.Len() != 3 {
		t.Errorf("after Delete, heap should have 3 elements, got %d", h.Len())
	}

	// Extract all remaining elements
	var extracted []int
	for h.Len() > 0 {
		extracted = append(extracted, h.ExtractMin())
	}

	expected := []int{1, 5, 7}
	if !slices.Equal(extracted, expected) {
		t.Errorf("extracted = %v, want %v", extracted, expected)
	}

	// Delete an already-deleted item should be safe
	item2.Delete()

	// Delete remaining items should be safe
	item1.Delete()
	item3.Delete()
	item4.Delete()
}

func TestItemFix(t *testing.T) {
	h := New[int]()

	// Insert some elements
	items := make([]Item, 5)
	items[0] = h.Insert(5)
	items[1] = h.Insert(3)
	items[2] = h.Insert(7)
	items[3] = h.Insert(1)
	items[4] = h.Insert(9)

	// Build the heap to establish invariant
	h.Build()

	// Modify the value at items[3] (currently 1) by accessing the internal data
	// In a real scenario, the user would modify their own data structure
	// For this test, we need to access the internal representation
	// Let's change the value and then call Fix

	// Since we can't directly modify through the Item, we'll test that
	// Fix maintains the heap invariant by modifying internal state
	idx := *items[3].index
	h.impl.data[idx].value = 8
	items[3].Fix()

	// Extract all elements - should still be in sorted order
	var extracted []int
	for h.Len() > 0 {
		extracted = append(extracted, h.ExtractMin())
	}

	expected := []int{3, 5, 7, 8, 9}
	if !slices.Equal(extracted, expected) {
		t.Errorf("after Fix, extracted = %v, want %v", extracted, expected)
	}
}

func TestClear(t *testing.T) {
	h := New[int]()

	items := make([]Item, 3)
	items[0] = h.Insert(5)
	items[1] = h.Insert(3)
	items[2] = h.Insert(7)

	h.Clear()

	if h.Len() != 0 {
		t.Errorf("after Clear, len should be 0, got %d", h.Len())
	}

	// Operations on items from cleared heap should be safe
	items[0].Delete()
	items[1].Fix()
}

func TestAll(t *testing.T) {
	h := New[int]()

	values := []int{5, 2, 8, 1, 9}
	for _, v := range values {
		h.Insert(v)
	}

	// Collect all elements
	var collected []int
	for v := range h.All() {
		collected = append(collected, v)
	}

	if len(collected) != 5 {
		t.Errorf("All() yielded %d elements, want 5", len(collected))
	}

	// First element should be the minimum
	if collected[0] != 1 {
		t.Errorf("first element from All() = %d, want 1", collected[0])
	}

	// All original values should be present
	slices.Sort(collected)
	expected := []int{1, 2, 5, 8, 9}
	if !slices.Equal(collected, expected) {
		t.Errorf("All() values = %v, want %v", collected, expected)
	}
}

func TestAllEarlyBreak(t *testing.T) {
	h := New[int]()

	for i := 0; i < 10; i++ {
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
	h := New[int]()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Min() on empty heap should panic")
		}
	}()

	h.Min()
}

func TestPanicOnEmptyExtractMin(t *testing.T) {
	h := New[int]()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ExtractMin() on empty heap should panic")
		}
	}()

	h.ExtractMin()
}

func TestDelayedHeapification(t *testing.T) {
	h := New[int]()

	// Insert elements without calling Build
	h.Insert(5)
	h.Insert(3)
	h.Insert(7)
	h.Insert(1)

	// First call to Min should trigger heapification
	if min := h.Min(); min != 1 {
		t.Errorf("Min() = %d, want 1", min)
	}

	// Subsequent inserts should maintain heap invariant
	h.Insert(0)
	if min := h.Min(); min != 0 {
		t.Errorf("after insert, Min() = %d, want 0", min)
	}
}

func TestHeapWithStrings(t *testing.T) {
	h := New[string]()

	h.Insert("dog")
	h.Insert("cat")
	h.Insert("bird")
	h.Insert("ant")

	if min := h.Min(); min != "ant" {
		t.Errorf("Min() = %q, want %q", min, "ant")
	}

	var extracted []string
	for h.Len() > 0 {
		extracted = append(extracted, h.ExtractMin())
	}

	expected := []string{"ant", "bird", "cat", "dog"}
	if !slices.Equal(extracted, expected) {
		t.Errorf("extracted = %v, want %v", extracted, expected)
	}
}

func TestLargeHeap(t *testing.T) {
	h := New[int]()

	// Insert 1000 elements in reverse order
	for i := 1000; i > 0; i-- {
		h.Insert(i)
	}

	// Extract all and verify they come out sorted
	prev := 0
	for h.Len() > 0 {
		curr := h.ExtractMin()
		if curr <= prev {
			t.Errorf("extracted %d after %d, not in sorted order", curr, prev)
		}
		prev = curr
	}

	if prev != 1000 {
		t.Errorf("last extracted value = %d, want 1000", prev)
	}
}
