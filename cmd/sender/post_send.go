package sender

import (
	"aws-ses-sender-go/config"
	"aws-ses-sender-go/model"
	"log"
	"time"
)

const (
	bulkSize   = 1000
	bulkPeriod = 10 * time.Second
)

// ConsumePostSend updates messages after processing
func ConsumePostSend() {
	buffer := make([]result, 0, bulkSize)

	ticker := time.NewTicker(bulkPeriod)
	defer ticker.Stop()

	for {
		select {
		case r := <-resultChan:
			buffer = append(buffer, r)
			if len(buffer) >= bulkSize {
				flushBuffer(&buffer)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				flushBuffer(&buffer)
			}
		}
	}
}

func flushBuffer(buf *[]result) {
	if len(*buf) == 0 {
		return
	}
	db := config.GetDB()
	tx := db.Begin()
	if err := tx.Error; err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return
	}

	for _, m := range *buf {
		tx.Model(&model.Request{}).
			Where("id = ?", m.ID).
			Updates(model.Request{
				MessageId: m.MessageId,
				Status:    m.Status,
				Error:     m.Error,
			})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("failed to commit updates: %v", err)
	} else {
		log.Printf("bulk update success: %d messages updated", len(*buf))
	}

	// Clear the buffer
	*buf = (*buf)[:0]
}
