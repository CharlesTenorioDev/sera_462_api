package sqsservice

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/sera_backend/internal/config"
)

// SqsInterface defines the methods for interacting with SQS.
type SqsInterface interface {
	SendMessage(ctx context.Context, msg *Message) error
	ReceiveMessages(callback func(msg types.Message))
	CreateQueue(ctx context.Context, queueName string, isFifo bool) (string, error)
}

// Message represents the structure of a message to be sent.
type Message struct {
	Data        []byte
	ContentType string
}

// SqsPool is the struct that implements the SqsInterface.
type SqsPool struct {
	client *sqs.Client
	conf   *config.Config
}

// NewSqsService initializes a new SqsPool with the given SQS client.
func NewSqsService(client *sqs.Client, conf *config.Config) SqsInterface {
	return &SqsPool{client: client, conf: conf}
}

// CreateQueue creates a new SQS queue.
func (s *SqsPool) CreateQueue(ctx context.Context, queueName string, isFifo bool) (string, error) {
	attributes := map[string]string{}
	if isFifo {
		attributes["FifoQueue"] = "true"
		queueName += ".fifo"
	}

	queue, err := s.client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName:  aws.String(queueName),
		Attributes: attributes,
	})
	if err != nil {
		log.Printf("Couldn't create queue %v: %v\n", queueName, err)
		return "", err
	}

	return *queue.QueueUrl, nil
}

// SendMessage sends a message to the specified SQS queue.
func (s *SqsPool) SendMessage(ctx context.Context, msg *Message) error {
	_, err := s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &s.conf.AWS_CONFIG.SQS_QUEUE_URL,
		MessageBody: aws.String(string(msg.Data)),
	})
	if err != nil {
		log.Printf("Failed to send message: %v\n", err)
		return err
	}
	return nil
}

// ReceiveMessages receives messages from the specified SQS queue and processes them with the callback.
func (s *SqsPool) ReceiveMessages(callback func(msg types.Message)) {
	for {
		output, err := s.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            &s.conf.AWS_CONFIG.SQS_QUEUE_URL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			log.Printf("Failed to receive messages: %v\n", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, msg := range output.Messages {
			callback(msg)

			// Delete message after processing
			_, err := s.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      &s.conf.AWS_CONFIG.SQS_QUEUE_URL,
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Printf("Failed to delete message: %v\n", err)
			}
		}
	}
}
