package models

import "time"

type Notification struct {
	ID        string                 `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Channel   string                 `gorm:"type:varchar(50);not null" json:"channel"` // email, sms, push
	Message   string                 `gorm:"type:text;not null" json:"message"`
	Metadata  map[string]interface{} `gorm:"type:jsonb" json:"metadata"`
	Status    string                 `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, processing, completed
	CreatedAt time.Time              `gorm:"autoCreateTime" json:"created_at"`
}

type NotificationRecipient struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	NotificationID string    `gorm:"type:uuid;not null;index" json:"notification_id"`
	Recipient      string    `gorm:"type:varchar(255);not null" json:"recipient"`
	Status         string    `gorm:"type:varchar(20);default:'queued'" json:"status"` // queued, sent, failed, retrying
	AttemptCount   int       `gorm:"default:0" json:"attempt_count"`
	LastError      string    `gorm:"type:text" json:"last_error,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type ProviderLog struct {
	ID                      string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	NotificationRecipientID string    `gorm:"type:uuid;not null;index" json:"notification_recipient_id"`
	Provider                string    `gorm:"type:varchar(50);not null" json:"provider"` // resend, twilio, fcm
	ResponseCode            string    `gorm:"type:varchar(50)" json:"response_code"`
	ResponseBody            string    `gorm:"type:text" json:"response_body"`
	Status                  string    `gorm:"type:varchar(20);not null" json:"status"` // success/failure
	CreatedAt               time.Time `gorm:"autoCreateTime" json:"created_at"`
}
