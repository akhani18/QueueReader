package poller

import (
	"log"

	"github.com/akhani18/QueueReader/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (p *Poller) poll(q *queue.Queue) {
	// Receive Messages.
	result, err := q.SQSClient.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            q.URL,
		MaxNumberOfMessages: aws.Int64(10),                // Batch receive upto 10 msgs
		WaitTimeSeconds:     aws.Int64(q.WaitTimeSeconds), // Long polling for 20 secs
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

	// Process Messages
	log.Printf("Poller#%d: Received %d messages.\n", p.pollerID, len(result.Messages))
	p.utilizationChan <- len(result.Messages)

	//log.Printf("Poller#%d: Received %d messages.\n", p.pollerID, 10)
	//p.utilizationChan <- 10

	for _, msg := range result.Messages {
		go p.process(msg, q)
	}
}

// Batch Delete
func (p *Poller) process(m *sqs.Message, q *queue.Queue) {
	// Process
	/*log.Printf("Poller#%d: Message Id: %s\n", p.pollerID, *m.MessageId)
	log.Printf("Poller#%d: Message ReceiveCount: %s\n", p.pollerID, *m.Attributes["ApproximateReceiveCount"])
	log.Printf("Poller#%d: Message Payload: %s\n", p.pollerID, *m.Body)*/

	// Delete
	_, err := q.SQSClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      q.URL,
		ReceiptHandle: m.ReceiptHandle,
	})

	if err != nil {
		log.Printf("Poller#%d: Error: %s", p.pollerID, err.Error())
		return
	}

	//log.Printf("Poller#%d: Successfully deleted message Id: %s\n", p.pollerID, *m.MessageId)
}
