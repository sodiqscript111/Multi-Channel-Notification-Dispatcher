package handlers

import (
	"NotifyX/db"
	"NotifyX/kafka"
	"NotifyX/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateNotificationRequest struct {
	Channel    string                 `json:"channel" binding:"required"` // "email", "sms", "push"
	Message    string                 `json:"message" binding:"required"`
	Recipients []string               `json:"recipients" binding:"required,min=1"`
	Metadata   map[string]interface{} `json:"metadata"` // optional
}

func CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest

	// 1) Validate Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 2) Create base notification entry
	notification := models.Notification{
		Channel:  req.Channel,
		Message:  req.Message,
		Metadata: req.Metadata,
		Status:   "pending",
	}

	// 3) Start database transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open transaction"})
		return
	}

	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification: " + err.Error()})
		return
	}

	// 4) Insert recipients + queue jobs
	queuedCount := 0
	for _, recipient := range req.Recipients {
		nr := models.NotificationRecipient{
			NotificationID: notification.ID,
			Recipient:      recipient,
			Status:         "queued",
			AttemptCount:   0,
		}

		if err := tx.Create(&nr).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipient: " + err.Error()})
			return
		}

		// ✅ Queue task (will implement enqueueRecipientJob next step)
		if err := kafka.PublishNotificationEvent(notification.ID, nr.ID); err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to send event to Kafka"})
			return
		}

		queuedCount++
	}

	// 5) Commit if everything succeeded
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// ✅ Response
	c.JSON(http.StatusCreated, gin.H{
		"notification_id":   notification.ID,
		"queued_recipients": queuedCount,
		"status":            "queued",
	})
}
