package median

import (
	"context"
)

func Process(ctx context.Context, rm RunningMedian, in <-chan int, out chan<- float64) error {
	if out != nil {
		defer close(out)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-in:
			if !ok {
				return nil
			}
			rm.Add(v)

			if out != nil {
				median := rm.Median()
				select {
				case out <- median:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}
