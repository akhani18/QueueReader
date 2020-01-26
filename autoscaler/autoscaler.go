package autoscaler

import (
	"fmt"
	"time"

	"github.com/akhani18/QueueReader/poller"
	"github.com/akhani18/QueueReader/queue"
)

const (
	MAX = 10
)

type Autoscaler struct {
	ps   []*poller.Poller
	uc   chan int
	minP int
	maxP int
	sqs  *queue.Queue
}

// New
func New(q *queue.Queue, minPollers int, maxPollers int) *Autoscaler {
	utilChan := make(chan int, 100)
	pollers := []*poller.Poller{poller.New(1, utilChan)}

	a := &Autoscaler{
		ps:   pollers,
		uc:   utilChan,
		minP: minPollers,
		maxP: maxPollers,
		sqs:  q,
	}

	return a
}

// Start
func (a *Autoscaler) Start() {
	// start monitor in a goroutine
	go a.monitor()

	// Start first poller
	go a.ps[0].Start(a.sqs)
}

// scaleUp by 1 poller
func (a *Autoscaler) scaleUp() {
	numPollers := len(a.ps)

	if numPollers < a.maxP {
		newPoller := poller.New(numPollers+1, a.uc)
		a.ps = append(a.ps, newPoller)

		go newPoller.Start(a.sqs)
	} else {
		fmt.Println("Can't scaleup, reached max num pollers")
	}
}

// scaleDown
func (a *Autoscaler) scaleDown() {
	numPollers := len(a.ps)

	if numPollers > a.minP {
		delPoller := a.ps[numPollers-1]
		delPoller.Stop()

		a.ps = a.ps[:numPollers-1]
	} else {
		fmt.Println("Can't scaledown, reached min num pollers")
	}
}

// monitor
func (a *Autoscaler) monitor() {
	ticker := time.NewTicker(60 * time.Second)
	sum := 0
	count := 0

	for {
		select {
		case t := <-ticker.C:
			fmt.Println("Time for autoscaling: ", t)
			if count != 0 {
				avg := float64(sum) / float64(count)

				if avg <= 0.75*MAX {
					fmt.Println("Performing scale down at ", t)
					a.scaleDown()
				} else if avg >= MAX {
					fmt.Println("Performing scale up at ", t)
					a.scaleUp()
				} else {
					fmt.Println("No autoscaling required")
				}
				sum = 0
				count = 0
			}

		case numMsg := <-a.uc:
			sum += numMsg
			count++
			//fmt.Printf("Accumulated messages %d, for count %d\n", sum, count)
		}

	}
}
