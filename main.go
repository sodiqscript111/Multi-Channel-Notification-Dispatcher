package main

import (
	"NotifyX/db"
	"NotifyX/handlers"
	"NotifyX/kafka"
	"NotifyX/server"
	"NotifyX/services"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db.ConnectDb()
	sqlDB, _ := db.DB.DB()

	kafka.InitKafkaProducer("localhost:9092")

	go kafka.StartNotificationConsumer(ctx)

	router := gin.Default()
	server.RegisterSystemRoutes(router)

	svc := services.NewNotificationService()
	router.POST("/notifications", handlers.CreateNotificationHandler(svc))
	server.MarkAsReady()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Println(" HTTP Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server crashed: %v", err)
		}
	}()

	// BLOCK until shutdown signal is received
	<-ctx.Done()
	log.Println("\n⚠️ Shutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Forced HTTP shutdown:", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}

	// Close Kafka producer
	if kafka.KafkaWriter != nil {
		log.Println("Closing Kafka producer...")
		_ = kafka.KafkaWriter.Close()
	}

	// Close DB connection
	log.Println("Closing database connection...")
	_ = sqlDB.Close()

	log.Println("Shutdown complete. Exiting cleanly.")
}
