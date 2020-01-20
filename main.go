package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/akhani18/QueueReader/queue"
)

func main() {
	qName := flag.String("n", "", "Name of the queue used for communication.")
	region := flag.String("r", "us-west-2", "AWS region where the queue exists.")
	longPollDurationSecs := flag.Int64("l", 20, "Long poll duration for SQS pollers (0 - 20 seconds).")
	numPollers := flag.Int("p", 1, "Number of concurrent pollers (1 - 20).")

	flag.Parse()

	if *qName == "" {
		flag.PrintDefaults()
		log.Println("Queue name is required.")
		os.Exit(1)
	}

	if *longPollDurationSecs < 0 || *longPollDurationSecs > 20 {
		flag.PrintDefaults()
		log.Println("Long poll duration should be between 0 and 20 seconds.")
		os.Exit(1)
	}

	if *numPollers < 1 || *numPollers > 20 {
		flag.PrintDefaults()
		log.Println("Number of concurrent pollers should be between 1 and 20.")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(*numPollers)

	for i := 1; i <= *numPollers; i++ {
		poller := queue.NewPoller(i, qName, region, *longPollDurationSecs)

		go poller.Run()
	}

	wg.Wait()
}
