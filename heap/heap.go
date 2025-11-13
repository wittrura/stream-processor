package heap

type MinHeap []int

// Len implements heap.Interface.
func (m MinHeap) Len() int {
	return len(m)
}

// Less implements heap.Interface.
func (m MinHeap) Less(i int, j int) bool {
	return m[i] < m[j]
}

// Swap implements heap.Interface.
func (m MinHeap) Swap(i int, j int) {
	m[i], m[j] = m[j], m[i]
}

// Push implements heap.Interface.
func (m *MinHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*m = append(*m, x.(int))
}

// Pop implements heap.Interface.
func (m *MinHeap) Pop() any {
	old := *m
	n := len(old)
	x := old[n-1]
	*m = old[0 : n-1]
	return x
}

type MaxHeap []int

// Len implements heap.Interface.
func (m MaxHeap) Len() int {
	return len(m)
}

// Less implements heap.Interface.
func (m MaxHeap) Less(i int, j int) bool {
	return m[i] > m[j]
}

// Swap implements heap.Interface.
func (m MaxHeap) Swap(i int, j int) {
	m[i], m[j] = m[j], m[i]
}

// Push implements heap.Interface.
func (m *MaxHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*m = append(*m, x.(int))
}

// Pop implements heap.Interface.
func (m *MaxHeap) Pop() any {
	old := *m
	n := len(old)
	x := old[n-1]
	*m = old[0 : n-1]
	return x
}
