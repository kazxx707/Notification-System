package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kavya/noti/handlers"
	"github.com/kavya/noti/repository"
	"github.com/kavya/noti/service"
	"github.com/kavya/noti/storage"
)

func main() {
	// Get database connection string from environment
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost/notifications?sslmode=disable"
		log.Printf("Using default DATABASE_URL: %s", connStr)
	}

	// Initialize database connection
	db, err := storage.NewDB(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	subscriptionRepo := repository.NewSubscriptionRepo(db.DB)
	notificationRepo := repository.NewNotificationRepo(db.DB)

	// Initialize service
	jsonFilePath := "notifications.json"
	notificationService := service.NewNotificationService(
		subscriptionRepo,
		notificationRepo,
		db.DB,
		jsonFilePath,
	)

	// Initialize handlers
	subscribeHandler := handlers.NewSubscribeHandler(notificationService)
	restockHandler := handlers.NewRestockHandler(notificationService)

	// Setup HTTP routes
	http.HandleFunc("/subscribe", subscribeHandler.Handle)
	http.HandleFunc("/inventory/restock", restockHandler.Handle)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Subscribe: POST http://localhost:%s/subscribe", port)
	log.Printf("Restock:   POST http://localhost:%s/inventory/restock", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
