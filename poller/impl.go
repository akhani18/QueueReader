package poller

import (
	"log"

	"github.com/akhani18/QueueReader/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	MaxMsgsPerPoll            = 10
	WaitDurationSec           = 20
	PollingFrequencyPerSecond = 1
)

func (p *Poller) poll(q *queue.Queue) {
	// Receive Messages
	result, err := q.SQSClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            q.URL,
		MaxNumberOfMessages: aws.Int64(MaxMsgsPerPoll),
		WaitTimeSeconds:     aws.Int64(WaitDurationSec),
		MessageAttributeNames: aws.StringSlice([]string{
			"All",
		}),
		AttributeNames: aws.StringSlice([]string{
			"All",
		}),
	})

	if err != nil {
		log.Printf("Poller#%d: Error: %s", p.pollerID, err.Error())
		return
	}

	// Send number of messages receieved on the channel for autoscaler to consume.
	p.utilizationChan <- len(result.Messages)

	if len(result.Messages) == 0 {
		log.Printf("Poller#%d: Empty poll, no messages receieved.\n", p.pollerID)
		return
	}

	// Process Messages - Dummy processing
	log.Printf("Poller#%d: Processed %d messages.\n", p.pollerID, len(result.Messages))

	// Delete Messages
	for _, msg := range result.Messages {
		go p.delete(msg, q)
	}
}

func (p *Poller) delete(m *sqs.Message, q *queue.Queue) {
	_, err := q.SQSClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      q.URL,
		ReceiptHandle: m.ReceiptHandle,
	})

	if err != nil {
		log.Printf("Poller#%d: Error: %s", p.pollerID, err.Error())
		return
	}
}
