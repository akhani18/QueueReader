package poller

import (
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
	ticker := time.NewTicker(PollingFrequencyPerSecond * time.Second)
	defer ticker.Stop()
	log.Printf("Poller#%d: Polling queue %s with a %d sec long-poll and %d sec period.\n", p.pollerID, *q.URL, WaitDurationSec, PollingFrequencyPerSecond)

	for {
		select {
		case <-ticker.C:
			p.poll(q)

		case <-p.stop:
			log.Printf("Poller#%d: Stopping...\n", p.pollerID)
			return
		}
	}
}

func (p *Poller) Stop() {
	p.stop <- true
}
