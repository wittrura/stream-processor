package median_test

import (
	"math"
	"runtime"
	"sync"
	"testing"
	"time"

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

func TestRunningMedian_ConcurrentAdds_CountMatches(t *testing.T) {
	m := NewMedianCalulator()

	numWorkers := runtime.NumCPU() * 4
	numPerWorker := 2000

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			defer wg.Done()
			base := workerID * numPerWorker
			for i := 0; i < numPerWorker; i++ {
				// Values are unique per (workerID, i) pair just
				// to mix things up a bit.
				v := base + i
				m.Add(v)
			}
		}(w)
	}

	wg.Wait()

	wantCount := numWorkers * numPerWorker
	gotCount := m.Count()
	if gotCount != wantCount {
		t.Fatalf("concurrent adds: Count()=%d, want=%d", gotCount, wantCount)
	}

	// Just sanity-check the final median is finite.
	median := m.Median()
	if math.IsNaN(median) || math.IsInf(median, 0) {
		t.Fatalf("concurrent adds: got invalid median %v", median)
	}
}

func TestRunningMedian_ConcurrentAddAndMedian_NoPanic(t *testing.T) {
	m := NewMedianCalulator()

	numAdders := runtime.NumCPU() * 2
	numPerAdder := 3000

	var wg sync.WaitGroup
	wg.Add(numAdders)

	stopReaders := make(chan struct{})

	// Reader goroutines: repeatedly call Median() and Count()
	// while writers are active. We don't assert much here
	// beyond "no panics / no data races" (checked via -race).
	numReaders := runtime.NumCPU() * 2
	for r := 0; r < numReaders; r++ {
		go func() {
			for {
				select {
				case <-stopReaders:
					return
				default:
					_ = m.Count()
					_ = m.Median()
					// Yield to avoid completely hogging CPU.
					runtime.Gosched()
				}
			}
		}()
	}

	for w := 0; w < numAdders; w++ {
		go func(workerID int) {
			defer wg.Done()
			base := workerID * numPerAdder
			for i := 0; i < numPerAdder; i++ {
				v := base - i // mix signs a bit
				m.Add(v)
			}
		}(w)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// writers finished
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for concurrent writers to finish")
	}

	close(stopReaders)

	// After everything settles, Count should equal numAdders * numPerAdder.
	wantCount := numAdders * numPerAdder
	gotCount := m.Count()
	if gotCount != wantCount {
		t.Fatalf("concurrent add+median: Count()=%d, want=%d", gotCount, wantCount)
	}

	median := m.Median()
	if math.IsNaN(median) || math.IsInf(median, 0) {
		t.Fatalf("concurrent add+median: got invalid median %v", median)
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
