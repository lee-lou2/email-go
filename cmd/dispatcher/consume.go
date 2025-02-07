package dispatcher

import (
	"aws-ses-sender-go/cmd/sender"
	"aws-ses-sender-go/pkg/aws"
	"context"
	"encoding/json"
	"log"
	"time"
)

func Run() {
	ctx := context.Background()
	sqsClient, err := aws.NewSQSClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create SQS client: %v", err)
	}

	queueUrl, err := sqsClient.GetOrCreateQueue(ctx)
	if err != nil {
		log.Fatalf("Failed to create (or verify) the queue: %v", err)
	}
	log.Printf("Queue URL: %s", *queueUrl)

	// Loop to consume (ReceiveMessage) messages
	for {
		messages, err := sqsClient.ReceiveMessages(ctx, queueUrl, 1, 10)
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
						TopicId string `json:"topicId"`
						Email   string `json:"email"`
						Subject string `json:"subject"`
						Content string `json:"content"`
					} `json:"messages"`
				}

				if m.Body != nil {
					if err := json.Unmarshal([]byte(*m.Body), &reqBody); err != nil {
						log.Printf("Failed to parse JSON message: %v", err)
					} else {
						// Request the sender to send the email
						for _, message := range reqBody.Messages {
							ctx := context.Background()
							sender.Request(message.TopicId, message.Email, message.Subject, message.Content, ctx)
						}
					}
				} else {
					log.Printf("Message body is nil")
				}

				// Delete the message after processing
				if err := sqsClient.DeleteMessage(ctx, queueUrl, m.ReceiptHandle); err != nil {
					log.Printf("Failed to delete message: %v", err)
				}
			}
		} else {
			// Wait for 3 seconds if no messages are present
			time.Sleep(3 * time.Second)
		}
	}
}
