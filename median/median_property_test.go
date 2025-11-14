package median_test

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"slices"

	. "example.com/stream-processor/median"
)

// medianOracle returns the exact median of the given slice of ints.
// It does not mutate the input.
func medianOracle(xs []int) float64 {
	if len(xs) == 0 {
		panic("medianOracle called with empty slice")
	}

	cp := slices.Clone(xs)
	sort.Ints(cp)

	n := len(cp)
	mid := n / 2
	if n%2 == 1 {
		return float64(cp[mid])
	}
	return float64(cp[mid-1]+cp[mid]) / 2.0
}

func TestRunningMedian_Property_RandomSequences(t *testing.T) {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	const (
		numSequences  = 100 // number of random test cases
		maxSeqLength  = 200 // max length of each sequence
		valueRangeMin = -1000
		valueRangeMax = 1000
	)

	for caseIdx := range numSequences {
		m := NewMedianCalulator()
		var seen []int

		seqLen := r.Intn(maxSeqLength) + 1 // [1, maxSeqLength]

		for i := range seqLen {
			v := r.Intn(valueRangeMax-valueRangeMin+1) + valueRangeMin
			seen = append(seen, v)
			m.Add(v)

			want := medianOracle(seen)
			got := m.Median()

			if m.Count() != len(seen) {
				t.Fatalf("seed=%d case=%d i=%d: Count()=%d, want=%d",
					seed, caseIdx, i, m.Count(), len(seen))
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("seed=%d case=%d i=%d: panic while comparing medians: %v", seed, caseIdx, i, r)
					}
				}()
				almostEqual(t, got, want)
			}()
		}
	}
}

func TestRunningMedian_Property_AllSameValue(t *testing.T) {
	m := NewMedianCalulator()
	const (
		n     = 1000
		value = 7
	)

	for i := range n {
		m.Add(value)
		got := m.Median()
		almostEqual(t, got, float64(value))

		if m.Count() != i+1 {
			t.Fatalf("after %d inserts, Count()=%d, want=%d", i+1, m.Count(), i+1)
		}
	}
}

func TestRunningMedian_Property_AlternatingExtremes(t *testing.T) {
	m := NewMedianCalulator()
	const (
		n    = 501 // odd, to have a clear median at the end
		low  = -1_000_000
		high = 1_000_000
	)

	var seen []int
	for i := range n {
		var v int
		if i%2 == 0 {
			v = low
		} else {
			v = high
		}
		seen = append(seen, v)
		m.Add(v)

		want := medianOracle(seen)
		got := m.Median()
		almostEqual(t, got, want)
	}
}
