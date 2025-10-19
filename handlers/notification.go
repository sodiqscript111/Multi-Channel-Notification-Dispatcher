package handlers

import (
	"NotifyX/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateNotificationRequest struct {
	Channel    string                 `json:"channel" binding:"required"`
	Message    string                 `json:"message" binding:"required"`
	Recipients []string               `json:"recipients" binding:"required,min=1"`
	Metadata   map[string]interface{} `json:"metadata"`
}

func CreateNotificationHandler(svc services.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateNotificationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, queued, err := svc.CreateNotification(req.Channel, req.Message, req.Metadata, req.Recipients)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"notification_id":   id,
			"queued_recipients": queued,
			"status":            "queued",
		})
	}
}
