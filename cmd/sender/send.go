package sender

import (
	"context"
	"email-go/config"
	"email-go/model"
	"email-go/pkg/aws"
	"strconv"
	"time"
)

func ConsumeSend() {
	rateStr := config.GetEnv("EMAIL_RATE", "14")
	rate, _ := strconv.Atoi(rateStr)
	ctx := context.Background()
	client, err := aws.GetClient(ctx)
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
			go func(m *message) {
				// Add code for the open event at the end of the body
				serverHost := config.GetEnv("SERVER_HOST", "http://localhost:3000")
				content := m.Content
				content += `<img src="` + serverHost + `/v1/events/open/?messageId=` + strconv.Itoa(int(m.ID)) + `">`
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				msgId, err := aws.SendSESEmail(
					client,
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
