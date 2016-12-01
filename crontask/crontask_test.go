package crontask_test

import (
	"context"
	"testing"

	"time"

	"sync"

	"github.com/gonzaloserrano/go-nzalo/crontask/crontask"
	"github.com/stretchr/testify/assert"
)

func TestCronExecutesACounterTask(t *testing.T) {
	assert := assert.New(t)

	// arrange
	mutex := &sync.RWMutex{}
	count := 0
	counter := func(ctx context.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		count++
		println(count)
	}
	ctx, cancel := context.WithCancel(context.Background())

	// act
	delay := time.Millisecond
	go crontask.Cron(ctx, counter, delay)

	times := 10
	time.Sleep(delay * time.Duration(times))
	cancel()
	// wait for it... dary!
	<-ctx.Done()

	// assert
	mutex.RLock()
	defer mutex.Unlock()
	assert.Equal(count, times)
}
