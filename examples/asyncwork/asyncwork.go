package asyncwork

import "fmt"

// A Message holds information for an async work.
type Message struct {
	command   string
	aggregate int
}

func (msg Message) String() string {
	return fmt.Sprintf("<%s:%d>", msg.command, msg.aggregate)
}

// NewMessageCollection returns a new list of messages from a map of data.
func NewMessageCollection(data map[string]int) []Message {
	msgs := []Message{}
	for cmd, aid := range data {
		msgs = append(msgs, Message{command: cmd, aggregate: aid})
	}
	return msgs
}

// A Worker is a func that gets a message and does something with it.
type Worker func(Message)

// A Supervisor registers workers and handles messages to them.
type Supervisor interface {
	Register(Worker)
	Handle(Message)
}

// NewSupervisor returns a new DefaultSupervisor.
func NewSupervisor() *DefaultSupervisor {
	return &DefaultSupervisor{workers: []Worker{}}
}

// A DefaultSupervisor is a supervisor that holds a list of workers so when
// handling a work message fans-out to each of them.
type DefaultSupervisor struct {
	workers []Worker
}

// Register adds a worker to the workers list.
func (s *DefaultSupervisor) Register(worker Worker) {
	s.workers = append(s.workers, worker)
}

// Handle sends the message to each worker.
func (s *DefaultSupervisor) Handle(msg Message) {
	for _, worker := range s.workers {
		go worker(msg)
	}
}
