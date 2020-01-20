package queue

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewPoller(id int, qName *string, region *string, longPollDurationSec int64) *Poller {
	sess, _ := session.NewSession(&aws.Config{
		Region: region,
	})

	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: qName,
	})

	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}

	qURL := result.QueueUrl

	return &Poller{
		pollerID:        id,
		queueURL:        qURL,
		waitDurationSec: longPollDurationSec,
		sqsClient:       svc,
	}
}
