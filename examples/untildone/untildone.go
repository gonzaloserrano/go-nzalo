package untildone

import (
	"context"
	"time"
)

func Do(ctx context.Context, step time.Duration) bool {
	defer println("finished")
	var steps int
	for {
		select {
		case <-ctx.Done():
			return true
		default:
			println("do step", steps)
			steps++
			time.Sleep(step)
		}
	}
	return false
}
