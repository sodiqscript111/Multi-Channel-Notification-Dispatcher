package kafka

import (
	"NotifyX/db"
	"NotifyX/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func StartNotificationConsumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "notifications.send",
		GroupID: "notifyx-consumers",
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		var event NotificationEvent
		_ = json.Unmarshal(msg.Value, &event)

		var recipient models.NotificationRecipient
		if err := db.DB.First(&recipient, "id = ?", event.RecipientID).Error; err != nil {
			continue
		}

		var notification models.Notification
		if err := db.DB.First(&notification, "id = ?", event.NotificationID).Error; err != nil {
			continue
		}

		success := simulateSend(notification.Channel, recipient.Recipient, notification.Message)

		if !success {
			fallback := fallbackChannel(notification.Channel)
			log.Printf("Primary %s failed, trying  fallback: %s", notification.Channel, fallback)
			success = simulateSend(fallback, recipient.Recipient, notification.Message)
		}

		if success {
			db.DB.Model(&recipient).Updates(map[string]interface{}{
				"status":        "sent",
				"attempt_count": recipient.AttemptCount + 1,
				"updated_at":    time.Now(),
			})
		} else {
			db.DB.Model(&recipient).Updates(map[string]interface{}{
				"status":        "failed",
				"attempt_count": recipient.AttemptCount + 1,
				"last_error":    "All channels failed",
				"updated_at":    time.Now(),
			})
		}
	}
}

func simulateSend(channel, recipient, message string) bool {
	if channel == "email" {
		fmt.Println("SENDING EMAIL to:", recipient, "| message:", message)
		return false
	}
	if channel == "sms" {
		fmt.Println("SENDING SMS to:", recipient, "| message:", message)
		return true
	}
	return false
}

func fallbackChannel(current string) string {
	if current == "email" {
		return "sms"
	}
	return "email"
}
