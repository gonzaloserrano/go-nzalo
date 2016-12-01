package crontask

import (
	"context"
	"time"
)

// A Task is a function that receives a context and runs.
type Task func(context.Context)

// Cron is responsible for executing a task each delay time.
func Cron(ctx context.Context, task Task, delay time.Duration) {
	// execute inmediatly
	go task(ctx)

	// execute every tick
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			println("tick")
			go task(ctx)
		case <-ctx.Done():
			return
		}
	}
}
