package main

import (
	"flag"
	"log"
	"os"

	"github.com/akhani18/QueueReader/autoscaler"
	"github.com/akhani18/QueueReader/queue"
)

func main() {
	qName := flag.String("n", "", "Name of the queue used for communication.")
	region := flag.String("r", "us-west-2", "AWS region where the queue exists.")

	flag.Parse()

	if *qName == "" {
		flag.PrintDefaults()
		log.Println("Queue name is required.")
		os.Exit(1)
	}

	sqs := queue.New(qName, region)
	autos := autoscaler.New(sqs)

	autos.Start()
}
