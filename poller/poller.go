package poller

import (
	"fmt"
	"log"
	"time"

	"github.com/akhani18/QueueReader/queue"
)

type Poller struct {
	pollerID        int
	utilizationChan chan<- int
	stop            chan bool
}

func New(id int, utilChan chan<- int) *Poller {
	return &Poller{
		pollerID:        id,
		utilizationChan: utilChan,
		stop:            make(chan bool),
	}
}

func (p *Poller) Start(q *queue.Queue) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	log.Printf("Poller#%d: Polling queue %s with a 20 sec long-poll and 1 sec period.\n", p.pollerID, *q.URL)

	for {
		select {
		case <-ticker.C:
			p.poll(q)

		case <-p.stop:
			fmt.Printf("Poller#%d: Stopping...\n", p.pollerID)
			return
		}
	}
}

func (p *Poller) Stop() {
	p.stop <- true
}
