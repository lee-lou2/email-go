package sender

import (
	"context"
	"email-go/config"
	"email-go/model"
)

// Request sends an email
func Request(planId, to, subject, content string, ctx context.Context) {
	// Validate data
	if to == "" || subject == "" || content == "" {
		return
	}

	// Save to database
	db := config.GetDB()
	emailMessage := &model.EmailMessage{
		PlanID:  planId,
		To:      to,
		Subject: subject,
		Content: content,
		Status:  model.EmailMessageStatusCreated,
	}
	db.Create(emailMessage)
	id := emailMessage.ID

	// Deliver message
	reqChan <- message{
		ID:      id,
		To:      to,
		Subject: subject,
		Content: content,
		Ctx:     ctx,
	}
}
