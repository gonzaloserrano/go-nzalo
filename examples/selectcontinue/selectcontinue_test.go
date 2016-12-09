package selectcontinue

import (
	"context"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ok := do(ctx)
		if !ok {
			t.Error("not ok")
		}
	}()
	time.Sleep(time.Second)
	cancel()
	<-ctx.Done()
}
