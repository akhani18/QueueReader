package queue

/*
import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Poller struct {
	queueURL        *string
	waitDurationSec int64
	sqsClient       *sqs.SQS
	pollerID        int
}

func (p *Poller) Run() {
	log.Printf("Poller#%d: Polling queue %s with long-poll duration %d secs.\n", p.pollerID, *p.queueURL, p.waitDurationSec)

	for {
		// Receive Messages.
		result, err := p.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            p.queueURL,
			MaxNumberOfMessages: aws.Int64(10),                // Batch receive upto 10 msgs
			WaitTimeSeconds:     aws.Int64(p.waitDurationSec), // Long polling to reduce empty receives
			MessageAttributeNames: aws.StringSlice([]string{
				"All",
			}),
			AttributeNames: aws.StringSlice([]string{
				"All",
			}),
		})

		if err != nil {
			log.Printf("Poller#%d: Error: %s", p.pollerID, err.Error())
			continue
		}

		// Process Messages
		log.Printf("Poller#%d: Received %d messages.\n", p.pollerID, len(result.Messages))

		var wg sync.WaitGroup
		for _, msg := range result.Messages {
			wg.Add(1)
			go p.process(msg, &wg)
		}
		wg.Wait()
	}
}

func (p *Poller) process(m *sqs.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	// Process
	log.Printf("Poller#%d: Message Id: %s\n", p.pollerID, *m.MessageId)
	log.Printf("Poller#%d: Message ReceiveCount: %s\n", p.pollerID, *m.Attributes["ApproximateReceiveCount"])
	log.Printf("Poller#%d: Message Payload: %s\n", p.pollerID, *m.Body)

	// Delete
	_, err := p.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      p.queueURL,
		ReceiptHandle: m.ReceiptHandle,
	})

	if err != nil {
		log.Printf("Poller#%d: Error: %s", p.pollerID, err.Error())
		return
	}

	log.Printf("Poller#%d: Successfully deleted message Id: %s\n", p.pollerID, *m.MessageId)
}
*/
