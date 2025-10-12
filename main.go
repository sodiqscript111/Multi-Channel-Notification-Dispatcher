package main

import (
	"NotifyX/db"
	"NotifyX/handlers"
	"NotifyX/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDb()
	kafka.InitKafkaProducer("localhost:9092")
	go kafka.StartNotificationConsumer()

	router := gin.Default()
	router.POST("/notifications", handlers.CreateNotification)
	router.Run(":8080")

}
