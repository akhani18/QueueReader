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

	/*sess, _ := session.NewSession(&aws.Config{
		Region: region,
		//Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
	})

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: qName,
	})

	queueURL := result.QueueUrl

	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(*numPollers)

	for i := 1; i <= *numPollers; i++ {
		go func(i int) {
			log.Printf("Running poller # %d to poll queue %s with long poll duration %d secs.\n", i, *queueURL, *longPollDurationSecs)
			for {
				// Receive Messages.
				result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
					QueueUrl:            queueURL,
					MaxNumberOfMessages: aws.Int64(10),        // Batch receive upto 10 msgs
					WaitTimeSeconds:     longPollDurationSecs, // Long polling to reduce empty receives
					MessageAttributeNames: aws.StringSlice([]string{
						"All",
					}),
					AttributeNames: aws.StringSlice([]string{
						"All",
					}),
				})

				if err != nil {
					log.Println("Error: ", err)
					continue
				}

				// Process Messages
				log.Printf("Receieved %d messages\n.", len(result.Messages))

				//var w sync.WaitGroup
				for _, msg := range result.Messages {
					//w.Add(1)

					//go func(m *sqs.Message) {
					func(m *sqs.Message) {
						//defer wg.Done()

						// process
						log.Printf("Message Id: %s\n", *m.MessageId)
						log.Printf("Message ReceieveCount: %s\n", *m.Attributes["ApproximateReceiveCount"])
						log.Printf("Message Payload: %s\n", *m.Body)

						// Delete
						_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
							QueueUrl:      queueURL,
							ReceiptHandle: m.ReceiptHandle,
						})

						if err != nil {
							log.Println("Error: ", err)
						}
					}(msg)

					//w.Wait()
				}
			}

		}(i)
	}*/

	var wg sync.WaitGroup
	wg.Add(*numPollers)

	for i := 1; i <= *numPollers; i++ {
		poller := queue.NewPoller(i, qName, region, *longPollDurationSecs)

		go poller.Run()
	}

	wg.Wait()
}
