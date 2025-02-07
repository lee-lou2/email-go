package model

import (
	"aws-ses-sender-go/config"

	"gorm.io/gorm"
)

const (
	EmailMessageStatusCreated = iota // Creation complete
	EmailMessageStatusSent           // Sent
	EmailMessageStatusFailed         // Failed
	EmailMessageStatusStopped        // Stopped
)

type Request struct {
	gorm.Model
	TopicId   string `json:"topic_id" gorm:"index;not null"`
	MessageId string `json:"message_id" gorm:"index;null;type:varchar(255)"`
	To        string `json:"to" gorm:"not null;type:varchar(255)"`
	Subject   string `json:"subject" gorm:"not null;type:varchar(255)"`
	Content   string `json:"content" gorm:"not null;type:text"`
	Status    int    `json:"status" gorm:"default:0;not null;type:tinyint"`
	Error     string `json:"error" gorm:"null;type:varchar(255)"`
}

func (m *Request) TableName() string {
	return "email_requests"
}

type Result struct {
	gorm.Model
	RequestId uint    `json:"request_id" gorm:"index;not null"`
	Request   Request `json:"request" gorm:"foreignKey:RequestId;references:ID"`
	Status    string  `json:"status" gorm:"not null;type:varchar(50)"`
	Raw       string  `json:"raw" gorm:"null;type:json"`
}

func (m *Result) TableName() string {
	return "email_results"
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&Request{})
	_ = db.AutoMigrate(&Result{})
}
