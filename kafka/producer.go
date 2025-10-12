package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func InitKafkaProducer(brokerURL string) {
	KafkaWriter = &kafka.Writer{
		Addr:         kafka.TCP(brokerURL),
		Topic:        "notifications.send",
		RequiredAcks: kafka.RequireAll,
		Balancer:     &kafka.LeastBytes{},
	}
	log.Println("Kafka Producer connected to:", brokerURL)
}

type NotificationEvent struct {
	NotificationID string `json:"notification_id"`
	RecipientID    string `json:"recipient_id"`
}

func PublishNotificationEvent(notificationID, recipientID string) error {
	event := NotificationEvent{
		NotificationID: notificationID,
		RecipientID:    recipientID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(recipientID), // helps kafka partition by recipient for ordering
		Value: payload,
		Time:  time.Now(),
	}

	err = KafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Println(" Failed to publish Kafka message:", err)
		return err
	}

	log.Println("Published event to Kafka →", notificationID, recipientID)
	return nil
}
