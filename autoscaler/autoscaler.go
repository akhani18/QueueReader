package autoscaler

import (
	"sync"

	"github.com/akhani18/QueueReader/poller"
	"github.com/akhani18/QueueReader/queue"
)

type Autoscaler struct {
	workerPool      []*poller.Poller
	utilizationChan chan int
	sqs             *queue.Queue
}

func New(q *queue.Queue) *Autoscaler {
	utilChan := make(chan int, 100) // TODO: Do we need a buffered channel?
	pool := []*poller.Poller{poller.New(1, utilChan)}

	a := &Autoscaler{
		workerPool:      pool,
		utilizationChan: utilChan,
		sqs:             q,
	}

	return a
}

func (a *Autoscaler) Start() {
	var wg sync.WaitGroup
	wg.Add(1)

	// start monitor in a goroutine
	go a.monitor()

	// Start the first poller
	go a.workerPool[0].Start(a.sqs)

	// Block here indefinetely
	wg.Wait()
}
