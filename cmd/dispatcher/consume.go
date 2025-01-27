package dispatcher

import (
	"context"
	"email-go/cmd/sender"
	"email-go/config"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func parseJSONBody(body *string, target interface{}) error {
	if body == nil {
		return fmt.Errorf("body is nil")
	}
	return json.Unmarshal([]byte(*body), target)
}

func Run() {
	sqsClient, err := config.GetQueue()
	if err != nil {
		log.Fatalf("Failed to create SQS client: %v", err)
	}

	// Example queue name to be used
	queueName := "email-queue"

	// Create the queue if it does not already exist
	queueUrl, err := createQueueIfNotExists(sqsClient, queueName)
	if err != nil {
		log.Fatalf("Failed to create (or verify) the queue: %v", err)
	}
	log.Printf("Queue URL: %s", *queueUrl)

	// Loop to consume (ReceiveMessage) messages
	for {
		messages, err := receiveMessages(sqsClient, queueUrl, 1, 10)
		if err != nil {
			log.Printf("Failed to receive message: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if len(messages) > 0 {
			// Process messages immediately if present
			for _, m := range messages {
				var reqBody struct {
					Messages []struct {
						PlanId  string `json:"planId"`
						Email   string `json:"email"`
						Subject string `json:"subject"`
						Content string `json:"content"`
					} `json:"messages"`
				}
				if e := parseJSONBody(m.Body, &reqBody); e != nil {
					log.Printf("Failed to parse JSON message: %v", e)
				} else {
					// Request the sender to send the email
					for _, message := range reqBody.Messages {
						ctx := context.Background()
						sender.Request(message.PlanId, message.Email, message.Subject, message.Content, ctx)
					}
				}

				// Delete the message after processing
				if err := deleteMessage(sqsClient, queueUrl, m.ReceiptHandle); err != nil {
					log.Printf("Failed to delete message: %v", err)
				}
			}
		} else {
			// Wait for 3 seconds if no messages are present
			time.Sleep(3 * time.Second)
		}
	}
}

func createQueueIfNotExists(client *sqs.Client, name string) (*string, error) {
	ctx := context.TODO()

	// Check if the queue already exists
	listOut, err := client.ListQueues(ctx, &sqs.ListQueuesInput{})
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
	out, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: &name,
	})
	if err != nil {
		return nil, err
	}

	return out.QueueUrl, nil
}

func receiveMessages(client *sqs.Client, queueUrl *string, maxMessages int32, waitTimeSeconds int32) ([]types.Message, error) {
	ctx := context.TODO()
	out, err := client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
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

func deleteMessage(client *sqs.Client, queueUrl *string, receiptHandle *string) error {
	ctx := context.TODO()
	_, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      queueUrl,
		ReceiptHandle: receiptHandle,
	})
	return err
}
