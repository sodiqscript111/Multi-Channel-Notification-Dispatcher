package main

import (
	"NotifyX/db"
	"NotifyX/handlers"
	"NotifyX/kafka"
	"NotifyX/server"
	"NotifyX/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	server.RegisterSystemRoutes(router)
	db.ConnectDb()
	kafka.InitKafkaProducer("localhost:9092")
	svc := services.NewNotificationService()
	go kafka.StartNotificationConsumer()

	router.POST("/notifications", handlers.CreateNotificationHandler(svc))
	server.MarkAsReady()
	router.Run(":8080")

}
