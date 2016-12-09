package untildone

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const delay = time.Second

// test with 100ms more than the delay time, eg $ go test -timeout 1100ms
func TestUntilDone(t *testing.T) {
	assert := assert.New(t)

	// arrange
	ctx, cancel := context.WithCancel(context.Background())
	finished := make(chan struct{})

	// act
	go func() {
		ok := Do(ctx, delay/10)
		assert.True(ok)
		close(finished)
	}()

	time.Sleep(delay)
	cancel()
	<-ctx.Done()
	// if it hadn't work we would be stuck here and test would timeout
	<-finished
}
