package median

import (
	"container/heap"
	"math"

	h "example.com/stream-processor/heap"
)

type MedianCalulator struct {
	lower h.MaxHeap
	upper h.MinHeap
	n     int
}

func NewMedianCalulator() *MedianCalulator {
	lower := h.MaxHeap{}
	heap.Init(&lower)

	upper := h.MinHeap{}
	heap.Init(&upper)
	return &MedianCalulator{
		lower: lower,
		upper: upper,
		n:     0,
	}
}

func (m *MedianCalulator) Median() float64 {
	if m.n == 0 {
		return math.NaN()
	}

	if m.n%2 == 0 {
		return (float64(m.lower[0]) + float64(m.upper[0])) / 2
	}
	if m.lower.Len() > m.upper.Len() {
		return float64(m.lower[0])
	}
	return float64(m.upper[0])
}

func (m *MedianCalulator) Add(x int) {
	m.n++

	if m.lower.Len() == 0 {
		heap.Push(&m.lower, x)
		return
	}

	if x <= m.lower[0] {
		heap.Push(&m.lower, x)
	} else {
		heap.Push(&m.upper, x)
	}
	m.rebalance()
}

func (m *MedianCalulator) rebalance() {
	if m.lower.Len()-m.upper.Len() > 1 {
		for m.lower.Len()-m.upper.Len() > 1 {
			temp := heap.Pop(&m.lower)
			heap.Push(&m.upper, temp)
		}
	}

	if m.upper.Len()-m.lower.Len() > 1 {
		for m.upper.Len()-m.lower.Len() > 1 {
			temp := heap.Pop(&m.upper)
			heap.Push(&m.lower, temp)
		}
	}
}

func (m *MedianCalulator) Count() int {
	return m.n
}

var _ RunningMedian = (*MedianCalulator)(nil)

type RunningMedian interface {
	Add(x int)
	Median() float64 // O(1)
	Count() int
}
