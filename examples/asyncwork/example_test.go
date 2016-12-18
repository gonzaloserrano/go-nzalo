package asyncwork_test

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gonzaloserrano/go-nzalo/examples/asyncwork"
)

var r = rand.New(rand.NewSource(99))

var msgs = asyncwork.NewMessageCollection(map[string]int{
	"RegisterCustomer":      86346,
	"AddAddressForCustomer": 65236,
	"SendWelcomeEmail":      61531,
	"UpdateCustomerProfile": 1876122,
})

// This example creates a supervisor, registers one worker and handles
// some messages. Also prints information for human visualitzation.
func ExampleSupervisorWithSingleWorker() {
	sup := asyncwork.NewSupervisor()
	wg := &sync.WaitGroup{}
	sup.Register(func(msg asyncwork.Message) {
		fmt.Printf("Received message %s, starting heavy work\n", msg)
		defer fmt.Printf("Heavy work for message %s ended.\n", msg)
		time.Sleep(time.Duration(r.Intn(20)) * time.Millisecond)
		wg.Done()
	})

	wg.Add(len(msgs))
	for _, msg := range msgs {
		sup.Handle(msg)
		time.Sleep(time.Second)
	}
	wg.Wait()

	// Unordered output: Received message <SendWelcomeEmail:61531>, starting heavy work
	// Received message <UpdateCustomerProfile:1876122>, starting heavy work
	// Received message <RegisterCustomer:86346>, starting heavy work
	// Received message <AddAddressForCustomer:65236>, starting heavy work
	// Heavy work for message <SendWelcomeEmail:61531> ended.
	// Heavy work for message <UpdateCustomerProfile:1876122> ended.
	// Heavy work for message <RegisterCustomer:86346> ended.
	// Heavy work for message <AddAddressForCustomer:65236> ended.
}
