package services

import (
	"NotifyX/db"
	"NotifyX/kafka"
	"NotifyX/models"
)

type NotificationService interface {
	CreateNotification(channel, message string, metadata map[string]interface{}, recipients []string) (string, int, error)
}

type notificationService struct{}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (n *notificationService) CreateNotification(channel, message string, metadata map[string]interface{}, recipients []string) (string, int, error) {
	notification := models.Notification{
		Channel:  channel,
		Message:  message,
		Metadata: metadata,
		Status:   "pending",
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		return "", 0, tx.Error
	}

	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		return "", 0, err
	}

	queuedCount := 0
	for _, recipient := range recipients {
		nr := models.NotificationRecipient{
			NotificationID: notification.ID,
			Recipient:      recipient,
			Status:         "queued",
			AttemptCount:   0,
		}

		if err := tx.Create(&nr).Error; err != nil {
			tx.Rollback()
			return "", 0, err
		}

		if err := kafka.PublishNotificationEvent(notification.ID, nr.ID); err != nil {
			tx.Rollback()
			return "", 0, err
		}
		queuedCount++
	}

	if err := tx.Commit().Error; err != nil {
		return "", 0, err
	}

	return notification.ID, queuedCount, nil
}
