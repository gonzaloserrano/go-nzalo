package asyncwork_test

import (
	"sync"
	"testing"
	"time"

	"github.com/gonzaloserrano/go-nzalo/examples/asyncwork"
)

func TestSupervisorRegistersSomeWorkersThatFinishAfterHandlingMessages(t *testing.T) {
	sup := asyncwork.NewSupervisor()

	wg := &sync.WaitGroup{}
	numWorkers := r.Intn(50)
	worker := func(msg asyncwork.Message) {
		wg.Done()
	}
	for i := 0; i < numWorkers; i++ {
		sup.Register(worker)
	}
	// since the supervisor sends every message to every worker
	wg.Add(numWorkers * len(msgs))
	for _, msg := range msgs {
		sup.Handle(msg)
	}

	// assert that the wait is less than some small time since
	// go test default timeout is 10 mins
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Millisecond * 10):
		// if the close channel is not closed after this time means
		// the workers have not finished and the test must fail.
		t.Fail()
	}
}

// this is just to print the example output.
func TestExample(t *testing.T) {
	ExampleSupervisorWithSingleWorker()
}
