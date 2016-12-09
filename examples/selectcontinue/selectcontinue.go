package selectcontinue

import (
	"context"
	"time"
)

func do(ctx context.Context) bool {
	var count int
	for {
		println("for", count)
		select {
		case <-ctx.Done():
			return true
		default:
			count++
			time.Sleep(time.Millisecond * 100)
			continue
		}
		return false
	}
	return false
}
