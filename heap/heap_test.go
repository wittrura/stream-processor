package heap_test

import (
	"container/heap"
	"testing"

	. "example.com/stream-processor/heap"
)

func TestMinHeap(t *testing.T) {
	h := &MinHeap{2, 1, 5}
	heap.Init(h)

	if h.Len() != 3 {
		t.Errorf("expected length: 3, got %d", h.Len())
	}

	heap.Push(h, 3)
	expected := []int{1, 2, 3, 5}
	for i := 0; h.Len() > 0; i++ {
		v := heap.Pop(h)
		if v != expected[i] {
			t.Errorf("pop error, expected value: %d, got %d", expected[i], v)
		}
	}

	if h.Len() != 0 {
		t.Errorf("empty heap should return length 0")
	}
}

func TestMaxHeap(t *testing.T) {
	h := &MaxHeap{2, 1, 5}
	heap.Init(h)

	if h.Len() != 3 {
		t.Errorf("expected length: 3, got %d", h.Len())
	}

	heap.Push(h, 3)
	expected := []int{5, 3, 2, 1}
	for i := 0; h.Len() > 0; i++ {
		v := heap.Pop(h)
		if v != expected[i] {
			t.Errorf("pop error, expected value: %d, got %d", expected[i], v)
		}
	}

	if h.Len() != 0 {
		t.Errorf("empty heap should return length 0")
	}
}
