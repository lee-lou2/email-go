package sender

import (
	"aws-ses-sender-go/config"
	"aws-ses-sender-go/model"
	"aws-ses-sender-go/pkg/aws"
	"context"
	"strconv"
	"time"
)

// ConsumeSend consumes the email sending requests
func ConsumeSend() {
	rateStr := config.GetEnv("EMAIL_RATE", "14")
	rate, _ := strconv.Atoi(rateStr)
	ctx := context.Background()
	sesClient, err := aws.NewSESClient(ctx)
	if err != nil {
		panic(err)
	}

	for {
		<-time.After(time.Second / time.Duration(rate))
		select {
		case msg := <-reqChan:
			if msg.Ctx.Err() != nil {
				// Stop immediately if there is a context error
				resultChan <- result{
					MessageId: "",
					ID:        msg.ID,
					Status:    model.EmailMessageStatusStopped,
					Error:     msg.Ctx.Err().Error(),
				}
				continue
			}
			go func(m *request) {
				// Add code for the open event at the end of the body
				serverHost := config.GetEnv("SERVER_HOST", "http://localhost:3000")
				content := m.Content
				content += `<img src="` + serverHost + `/v1/events/open/?requestId=` + strconv.Itoa(int(m.ID)) + `">`
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				msgId, err := sesClient.SendEmail(
					ctx,
					&m.Subject,
					&content,
					&[]string{m.To},
				)
				if err != nil {
					// Sending failed
					resultChan <- result{
						MessageId: msgId,
						ID:        msg.ID,
						Status:    model.EmailMessageStatusFailed,
						Error:     err.Error(),
					}
					return
				}
				// Sending succeeded
				resultChan <- result{
					MessageId: msgId,
					ID:        msg.ID,
					Status:    model.EmailMessageStatusSent,
					Error:     "",
				}
			}(&msg)
		}
	}
}
