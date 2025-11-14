package median_test

import (
	"math"
	"testing"

	. "example.com/stream-processor/median"
)

func TestRunningMedian_EmptyReturnsNaN(t *testing.T) {
	m := NewMedianCalulator()

	if m.Count() != 0 {
		t.Fatalf("expected Count()=0 for new median, got %d", m.Count())
	}

	got := m.Median()
	if !math.IsNaN(got) {
		t.Fatalf("expected Median() to be NaN for empty stream, got %v", got)
	}
}

func TestRunningMedian_SingleElement(t *testing.T) {
	m := NewMedianCalulator()

	m.Add(42)

	if m.Count() != 1 {
		t.Fatalf("expected Count()=1, got %d", m.Count())
	}

	got := m.Median()
	almostEqual(t, got, 42.0)
}

func TestRunningMedian_OddCount_MedianIsMiddleElement(t *testing.T) {
	m := NewMedianCalulator()

	// Insert in unsorted order on purpose.
	input := []int{5, 1, 9, -2, 7}
	for _, v := range input {
		m.Add(v)
	}

	if m.Count() != len(input) {
		t.Fatalf("expected Count()=%d, got %d", len(input), m.Count())
	}

	// Sorted: -2, 1, 5, 7, 9 -> median = 5
	got := m.Median()
	almostEqual(t, got, 5.0)
}

func TestRunningMedian_EvenCount_MedianIsAverageOfTwoMiddles(t *testing.T) {
	m := NewMedianCalulator()

	input := []int{5, 1, 9, -2}
	for _, v := range input {
		m.Add(v)
	}

	if m.Count() != len(input) {
		t.Fatalf("expected Count()=%d, got %d", len(input), m.Count())
	}

	// Sorted: -2, 1, 5, 9 -> middle two are 1 and 5 => (1+5)/2 = 3
	got := m.Median()
	almostEqual(t, got, 3.0)
}

func TestRunningMedian_IncrementalMedians_MixedSequence(t *testing.T) {
	m := NewMedianCalulator()

	input := []int{5, 2, 10, -1, 6, 6}
	// After each Add, this is the expected median.
	// Sorted prefixes:
	// [5]                          -> 5
	// [2,5]                        -> (2+5)/2 = 3.5
	// [2,5,10]                     -> 5
	// [-1,2,5,10]                  -> (2+5)/2 = 3.5
	// [-1,2,5,6,10]                -> 5
	// [-1,2,5,6,6,10]              -> (5+6)/2 = 5.5
	wantMedians := []float64{5, 3.5, 5, 3.5, 5, 5.5}

	for i, v := range input {
		m.Add(v)
		got := m.Median()
		almostEqual(t, got, wantMedians[i])

		if m.Count() != i+1 {
			t.Fatalf("after %d inserts, expected Count()=%d, got %d", i+1, i+1, m.Count())
		}
	}
}

func TestRunningMedian_NegativesAndDuplicates(t *testing.T) {
	m := NewMedianCalulator()

	input := []int{-5, -5, -1, -10, -1, -5}
	// Sorted: -10, -5, -5, -5, -1, -1
	// middle two: -5 and -5 -> median = -5
	for _, v := range input {
		m.Add(v)
	}

	if m.Count() != len(input) {
		t.Fatalf("expected Count()=%d, got %d", len(input), m.Count())
	}

	got := m.Median()
	almostEqual(t, got, -5.0)
}

func TestRunningMedian_MonotonicIncreasing(t *testing.T) {
	m := NewMedianCalulator()

	input := []int{1, 2, 3, 4, 5, 6}
	wantMedians := []float64{
		1,   // [1]
		1.5, // [1,2]
		2,   // [1,2,3]
		2.5, // [1,2,3,4]
		3,   // [1,2,3,4,5]
		3.5, // [1,2,3,4,5,6]
	}

	for i, v := range input {
		m.Add(v)
		got := m.Median()
		almostEqual(t, got, wantMedians[i])
	}
}

func TestRunningMedian_MonotonicDecreasing(t *testing.T) {
	m := NewMedianCalulator()

	input := []int{6, 5, 4, 3, 2, 1}
	wantMedians := []float64{
		6,   // [6]
		5.5, // [5,6]
		5,   // [4,5,6]
		4.5, // [3,4,5,6]
		4,   // [2,3,4,5,6]
		3.5, // [1,2,3,4,5,6]
	}

	for i, v := range input {
		m.Add(v)
		got := m.Median()
		almostEqual(t, got, wantMedians[i])
	}
}

func almostEqual(t *testing.T, got, want float64) {
	t.Helper()
	const eps = 1e-9
	if math.IsNaN(want) {
		if !math.IsNaN(got) {
			t.Fatalf("expected NaN, got %v", got)
		}
		return
	}
	if math.Abs(got-want) > eps {
		t.Fatalf("floats not equal: got=%v want=%v", got, want)
	}
}
