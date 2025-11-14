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

func TestMinHeap_PushPopAscendingOrder(t *testing.T) {
	h := &MinHeap{}
	heap.Init(h)

	input := []int{5, 1, 9, -2, 7, 7, 0}
	for _, v := range input {
		heap.Push(h, v)
	}

	var got []int
	for h.Len() > 0 {
		got = append(got, heap.Pop(h).(int))
	}

	want := []int{-2, 0, 1, 5, 7, 7, 9}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("minHeap order mismatch at index %d: got %v, want %v (full=%v)", i, got[i], want[i], got)
		}
	}
}

func TestMaxHeap_PushPopDescendingOrder(t *testing.T) {
	h := &MaxHeap{}
	heap.Init(h)

	input := []int{5, 1, 9, -2, 7, 7, 0}
	for _, v := range input {
		heap.Push(h, v)
	}

	var got []int
	for h.Len() > 0 {
		got = append(got, heap.Pop(h).(int))
	}

	want := []int{9, 7, 7, 5, 1, 0, -2}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("maxHeap order mismatch at index %d: got %v, want %v (full=%v)", i, got[i], want[i], got)
		}
	}
}

func TestMinHeap_LenAndPushPop(t *testing.T) {
	h := &MinHeap{}
	heap.Init(h)

	if h.Len() != 0 {
		t.Fatalf("expected empty heap length 0, got %d", h.Len())
	}

	heap.Push(h, 3)
	heap.Push(h, 1)
	heap.Push(h, 2)

	if h.Len() != 3 {
		t.Fatalf("expected length 3, got %d", h.Len())
	}

	if got := heap.Pop(h).(int); got != 1 {
		t.Fatalf("expected first popped element to be 1, got %d", got)
	}
	if got := heap.Pop(h).(int); got != 2 {
		t.Fatalf("expected second popped element to be 2, got %d", got)
	}
	if got := heap.Pop(h).(int); got != 3 {
		t.Fatalf("expected third popped element to be 3, got %d", got)
	}
}

func TestMaxHeap_LenAndPushPop(t *testing.T) {
	h := &MaxHeap{}
	heap.Init(h)

	if h.Len() != 0 {
		t.Fatalf("expected empty heap length 0, got %d", h.Len())
	}

	heap.Push(h, 3)
	heap.Push(h, 1)
	heap.Push(h, 2)

	if h.Len() != 3 {
		t.Fatalf("expected length 3, got %d", h.Len())
	}

	if got := heap.Pop(h).(int); got != 3 {
		t.Fatalf("expected first popped element to be 3, got %d", got)
	}
	if got := heap.Pop(h).(int); got != 2 {
		t.Fatalf("expected second popped element to be 2, got %d", got)
	}
	if got := heap.Pop(h).(int); got != 1 {
		t.Fatalf("expected third popped element to be 1, got %d", got)
	}
}
