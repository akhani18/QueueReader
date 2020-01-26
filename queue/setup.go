package queue

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Queue struct {
	URL             *string
	WaitTimeSeconds int64
	SQSClient       *sqs.SQS
}

func New(qName *string, region *string, longPollDurationSec int64) *Queue {
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

	return &Queue{
		URL:             qURL,
		WaitTimeSeconds: longPollDurationSec,
		SQSClient:       svc,
	}
}
