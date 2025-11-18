package median_test

import (
	"context"
	"math"
	"testing"
	"time"

	. "example.com/stream-processor/median"
)

func TestProcess_EmitsMediansAndClosesOut_OnInputClose(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := make(chan int)
	out := make(chan float64)

	rm := NewMedianCalulator()

	errCh := make(chan error, 1)
	go func() {
		errCh <- Process(ctx, rm, in, out)
	}()

	// Send a small, known sequence.
	go func() {
		seq := []int{1, 5, 2, 10}
		for _, v := range seq {
			in <- v
		}
		close(in)
	}()

	// Collect all medians until out is closed.
	gotMedians := collectWithTimeout(t, out, 2*time.Second)

	// Expected medians after each insert:
	// [1]           -> 1
	// [1,5]         -> (1+5)/2 = 3
	// [1,2,5]       -> 2
	// [1,2,5,10]    -> (2+5)/2 = 3.5
	wantMedians := []float64{1, 3, 2, 3.5}

	if len(gotMedians) != len(wantMedians) {
		t.Fatalf("unexpected number of medians: got=%d want=%d (values=%v)",
			len(gotMedians), len(wantMedians), gotMedians)
	}

	for i := range wantMedians {
		if math.Abs(gotMedians[i]-wantMedians[i]) > 1e-9 {
			t.Fatalf("median mismatch at index %d: got=%v want=%v (all=%v)",
				i, gotMedians[i], wantMedians[i], gotMedians)
		}
	}

	// Process should return cleanly (no context cancel).
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Process returned error on normal completion: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for Process to return on normal completion")
	}

	// Final median from rm should match last emitted median.
	finalMedian := rm.Median()
	if math.Abs(finalMedian-wantMedians[len(wantMedians)-1]) > 1e-9 {
		t.Fatalf("final Median() from rm=%v, want=%v", finalMedian, wantMedians[len(wantMedians)-1])
	}
}

func TestProcess_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan int)
	out := make(chan float64)
	rm := NewMedianCalulator()

	errCh := make(chan error, 1)
	go func() {
		errCh <- Process(ctx, rm, in, out)
	}()

	// Start feeding values slowly.
	go func() {
		defer close(in)
		for i := range 1000 {
			in <- i
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Let some values through, then cancel the context.
	time.Sleep(10 * time.Millisecond)
	cancel()

	// We don't care about all medians; just make sure we can drain 'out'
	// until it's closed and Process returns.
	_ = collectWithTimeout(t, out, 2*time.Second)

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatalf("expected Process to return non-nil error on context cancel, got nil")
		}
		if err != context.Canceled && err != context.DeadlineExceeded {
			t.Fatalf("expected context-related error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for Process to return after context cancel")
	}
}

func TestProcess_StopsOnContextCancel_WhileIdle(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := make(chan int) // no one will ever send on this
	out := make(chan float64)

	rm := NewMedianCalulator()
	errCh := make(chan error, 1)

	go func() {
		errCh <- Process(ctx, rm, in, out)
	}()

	// Give Process time to start and block on the outer select.
	time.Sleep(10 * time.Millisecond)

	// Now cancel the context while it's idle.
	cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatalf("expected non-nil error on context cancel, got nil")
		}
		if err != context.DeadlineExceeded {
			t.Fatalf("expected context-related error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for Process to return after context cancel while idle")
	}
}

func TestProcess_AllowsNilOutChannel(t *testing.T) {
	// When out == nil, Process should still consume the stream and exit cleanly.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := make(chan int)
	var out chan float64 // nil

	rm := NewMedianCalulator()
	errCh := make(chan error, 1)

	go func() {
		errCh <- Process(ctx, rm, in, out)
	}()

	go func() {
		for _, v := range []int{10, 20, 30, 40} {
			in <- v
		}
		close(in)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Process(out=nil) returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for Process(out=nil) to return")
	}

	// rm should have seen all values.
	if rm.Count() != 4 {
		t.Fatalf("expected rm.Count()=4 with out=nil, got %d", rm.Count())
	}
	median := rm.Median()
	// Values: [10,20,30,40] -> median is (20+30)/2 = 25
	if math.Abs(median-25) > 1e-9 {
		t.Fatalf("unexpected median with out=nil: got=%v want=%v", median, 25.0)
	}
}

func collectWithTimeout[T any](t *testing.T, ch <-chan T, timeout time.Duration) []T {
	t.Helper()
	var out []T
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, v)
		case <-timer.C:
			t.Fatalf("timed out waiting for channel to close")
		}
	}
}
