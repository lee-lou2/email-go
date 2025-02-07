package sender

import (
	"aws-ses-sender-go/config"
	"aws-ses-sender-go/model"
	"context"
)

// Request sends an email
func Request(topicId, to, subject, content string, ctx context.Context) {
	// Validate data
	if to == "" || subject == "" || content == "" {
		return
	}

	// Save to database
	db := config.GetDB()
	emailMessage := &model.Request{
		TopicId: topicId,
		To:      to,
		Subject: subject,
		Content: content,
		Status:  model.EmailMessageStatusCreated,
	}
	db.Create(emailMessage)
	id := emailMessage.ID

	// Deliver request
	reqChan <- request{
		ID:      id,
		To:      to,
		Subject: subject,
		Content: content,
		Ctx:     ctx,
	}
}
