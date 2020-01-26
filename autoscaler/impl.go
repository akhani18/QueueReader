package autoscaler

import (
	"log"
	"time"

	"github.com/akhani18/QueueReader/poller"
)

const (
	MinNumPollers            = 1
	MaxNumPollers            = 25
	AutoscalingPeriodSeconds = 60
	ScaleupThreshold         = 0.90
	ScaledownThreshold       = 0.75
)

// scaleUp creates a new poller, starts it and adds it to the worker pool.
// Only scales up till MaxNumPolllers.
func (a *Autoscaler) scaleUp() {
	numPollers := len(a.workerPool)

	if numPollers < MaxNumPollers {
		newPoller := poller.New(numPollers+1, a.utilizationChan)
		a.workerPool = append(a.workerPool, newPoller)

		go newPoller.Start(a.sqs)
	} else {
		log.Println("Autoscaler: Can't scaleup, reached max num pollers.")
	}
}

// scaleDown sends stop signal to the last poller in the pool and removes it from the pool.
// Only scales down till MinNumPollers
func (a *Autoscaler) scaleDown() {
	numPollers := len(a.workerPool)

	if numPollers > MinNumPollers {
		delPoller := a.workerPool[numPollers-1]
		delPoller.Stop()

		a.workerPool = a.workerPool[:numPollers-1]
	} else {
		log.Println("Autoscaler: Can't scaledown, reached min num pollers.")
	}
}

// monitor listens for number of messages received by every poll and aggregates them till the ticker
// ticks, which is when it calls scaleUp, scaleDown or ignores depending on the average number of messages
// received per poll.
func (a *Autoscaler) monitor() {
	log.Printf("Autoscalar: Starting autoscaling monitor to run every %d seconds.\n", AutoscalingPeriodSeconds)
	ticker := time.NewTicker(AutoscalingPeriodSeconds * time.Second)
	totalMessages := 0
	numReceives := 0

	for {
		select {
		case <-ticker.C:
			if numReceives != 0 {
				avgMessagesPerPoll := float64(totalMessages) / float64(numReceives)

				if avgMessagesPerPoll < ScaledownThreshold*poller.MaxMsgsPerPoll {
					log.Println("Autoscaler: Attempting to perform scale down.")
					a.scaleDown()
				} else if avgMessagesPerPoll > ScaleupThreshold*poller.MaxMsgsPerPoll {
					log.Println("Autoscaler: Attempting to perform scale up.")
					a.scaleUp()
				} else {
					log.Println("Autoscaler: No autoscaling required.")
				}

				// Reset counters
				totalMessages = 0
				numReceives = 0
			}

		case numMsg := <-a.utilizationChan:
			totalMessages += numMsg
			numReceives++
		}
	}
}
