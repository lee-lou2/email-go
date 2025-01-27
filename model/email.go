package model

import (
	"email-go/config"

	"gorm.io/gorm"
)

const (
	EmailMessageStatusCreated = iota // Creation complete
	EmailMessageStatusSent           // Sent
	EmailMessageStatusFailed         // Failed
	EmailMessageStatusStopped        // Stopped
)

type EmailMessage struct {
	gorm.Model
	PlanID    string `json:"plan_id" gorm:"index;not null"`
	MessageId string `json:"message_id" gorm:"index;null;type:varchar(255)"`
	To        string `json:"to" gorm:"not null;type:varchar(255)"`
	Subject   string `json:"subject" gorm:"not null;type:varchar(255)"`
	Content   string `json:"content" gorm:"not null;type:text"`
	Status    int    `json:"status" gorm:"default:0;not null;type:tinyint"`
	Error     string `json:"error" gorm:"null;type:varchar(255)"`
}

func (m *EmailMessage) TableName() string {
	return "email_messages"
}

type EmailMessageResult struct {
	gorm.Model
	MessageId uint         `json:"message_id" gorm:"index;not null"`
	Message   EmailMessage `json:"message" gorm:"foreignKey:MessageId;references:ID"`
	Status    string       `json:"status" gorm:"not null;type:varchar(50)"`
	Raw       string       `json:"raw" gorm:"null;type:json"`
}

func (m *EmailMessageResult) TableName() string {
	return "email_message_results"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&EmailMessage{})
	_ = db.AutoMigrate(&EmailMessageResult{})
}
