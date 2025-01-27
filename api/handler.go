package api

import (
	"bytes"
	"email-go/config"
	"email-go/model"
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"image"
	"image/color"
	"image/png"
	"log"
	"strconv"
	"time"
)

// createOpenEventHandler Open Event Handler
// Attach an image script to the email and assume it has been read when the image is accessed
func createOpenEventHandler(c fiber.Ctx) error {
	msgId := c.Query("messageId")
	if msgId != "" {
		// Consider email as opened and create data
		db := config.GetDB()
		var message model.EmailMessageResult
		msgIdInt, _ := strconv.Atoi(msgId)
		message.MessageId = uint(msgIdInt)
		message.Status = "Open"
		_ = db.Create(&message).Error
	}

	// Return a blank image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 0, G: 0, B: 0, A: 0})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	c.Set("Content-Type", "image/png")
	return c.Send(buf.Bytes())
}

// createResultEventHandler Result Event Handler
// Handler that receives AWS SES results
func createResultEventHandler(c fiber.Ctx) error {
	var reqBody struct {
		Type         string `json:"Type"`
		Message      string `json:"Message"`
		SubscribeURL string `json:"SubscribeURL"`
	}
	if err := c.Bind().JSON(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if reqBody.Type == "SubscriptionConfirmation" {
		log.Println(reqBody.SubscribeURL)
		return c.JSON(fiber.Map{})
	} else if reqBody.Type != "Notification" {
		return c.JSON(fiber.Map{})
	}

	// Retrieve message
	var bodyMessage struct {
		EventType string `json:"eventType"`
		Mail      struct {
			MessageId string `json:"messageId"`
		} `json:"mail"`
	}
	if err := json.Unmarshal([]byte(reqBody.Message), &bodyMessage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Retrieve message
	db := config.GetDB()
	var message model.EmailMessage
	if err := db.Where("message_id = ?", bodyMessage.Mail.MessageId).First(&message).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save result
	result := model.EmailMessageResult{
		MessageId: message.ID,
		Status:    bodyMessage.EventType,
		Raw:       reqBody.Message,
	}
	_ = db.Create(&result).Error
	return c.JSON(fiber.Map{})
}

// getResultCountHandler Retrieve email delivery results as counts
func getResultCountHandler(c fiber.Ctx) error {
	// Get the planId from the URL parameters
	planId := c.Params("planId")
	if planId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "planId is required"})
	}

	db := config.GetDB()

	// Get the number of emails by status
	var results []struct {
		Status int
		Count  int
	}

	err := db.Model(&model.EmailMessage{}).
		Select("status, COUNT(*) as count").
		Where("plan_id = ?", planId).
		Group("status").
		Scan(&results).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the result
	var counts struct {
		Total   int `json:"total"`
		Created int `json:"created"`
		Sent    int `json:"sent"`
		Failed  int `json:"failed"`
		Stopped int `json:"stopped"`
	}

	for _, r := range results {
		switch r.Status {
		case model.EmailMessageStatusCreated:
			counts.Created = r.Count
		case model.EmailMessageStatusSent:
			counts.Sent = r.Count
		case model.EmailMessageStatusFailed:
			counts.Failed = r.Count
		case model.EmailMessageStatusStopped:
			counts.Stopped = r.Count
		}
		counts.Total += r.Count
	}
	return c.JSON(counts)
}

// getSentCountHandler Retrieve the number of emails sent within 24 hours
func getSentCountHandler(c fiber.Ctx) error {
	// Receive hours as a query string (default 24)
	hours, err := strconv.Atoi(c.Query("hours", "24"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Calculate the time hours before the current time
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	// Get the number of emails after startTime from the DB
	db := config.GetDB()
	var count int64
	if err := db.Model(&model.EmailMessage{}).
		Where("created_at > ?", startTime).
		Where("status = ?", model.EmailMessageStatusSent).
		Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return the result
	return c.JSON(fiber.Map{"count": count})
}
