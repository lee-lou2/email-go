package aws

import (
	"aws-ses-sender-go/config"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQS is a wrapper around the AWS SQS client
type SQS struct {
	Client *sqs.Client
}

// NewSQSClient creates a new SQS client
func NewSQSClient(ctx context.Context) (*SQS, error) {
	AccessKeyId := config.GetEnv("AWS_ACCESS_KEY_ID")
	SecretAccessKey := config.GetEnv("AWS_SECRET_ACCESS_KEY")
	cfg, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithRegion("ap-northeast-2"),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(AccessKeyId, SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &SQS{
		Client: sqs.NewFromConfig(cfg),
	}, nil
}

// GetOrCreateQueue gets or creates an SQS queue
func (s *SQS) GetOrCreateQueue(ctx context.Context) (*string, error) {
	name := config.GetEnv("AWS_SQS_DEFAULT_QUEUE_NAME", "email-queue")

	// Check if the queue already exists
	listOut, err := s.Client.ListQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		return nil, err
	}

	var existingUrl *string
	for _, url := range listOut.QueueUrls {
		if len(url) > 0 {
			if len(url) >= len(name) && url[len(url)-len(name):] == name {
				existingUrl = &url
				break
			}
		}
	}

	if existingUrl != nil {
		return existingUrl, nil
	}

	// If the queue does not exist, create a new one
	out, err := s.Client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: &name,
	})
	if err != nil {
		return nil, err
	}

	return out.QueueUrl, nil
}

// SendMessage sends a message to the SQS queue
func (s *SQS) SendMessage(ctx context.Context, queueUrl string, body string) (*sqs.SendMessageOutput, error) {
	input := &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(body),
	}
	return s.Client.SendMessage(ctx, input)
}

// ReceiveMessages receives messages from the SQS queue
func (s *SQS) ReceiveMessages(ctx context.Context, queueUrl *string, maxMessages int32, waitTimeSeconds int32) ([]types.Message, error) {
	out, err := s.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTimeSeconds,
		VisibilityTimeout:   30,
	})
	if err != nil {
		return nil, err
	}
	return out.Messages, nil
}

// DeleteMessage deletes a message from the SQS queue
func (s *SQS) DeleteMessage(ctx context.Context, queueUrl *string, receiptHandle *string) error {
	_, err := s.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      queueUrl,
		ReceiptHandle: receiptHandle,
	})
	return err
}
