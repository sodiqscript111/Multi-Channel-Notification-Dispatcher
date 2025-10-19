package main

import (
	"NotifyX/db"
	"NotifyX/handlers"
	"NotifyX/kafka"
	"NotifyX/services"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDb()
	kafka.InitKafkaProducer("localhost:9092")
	svc := services.NewNotificationService()
	go kafka.StartNotificationConsumer()

	router := gin.Default()
	router.POST("/notifications", handlers.CreateNotificationHandler(svc))
	router.Run(":8080")

}
