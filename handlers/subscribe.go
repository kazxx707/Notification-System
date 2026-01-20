package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kavya/noti/service"
)

// SubscribeHandler handles subscription requests
type SubscribeHandler struct {
	notificationService *service.NotificationService
}

// NewSubscribeHandler creates a new subscribe handler
func NewSubscribeHandler(notificationService *service.NotificationService) *SubscribeHandler {
	return &SubscribeHandler{
		notificationService: notificationService,
	}
}

// SubscribeRequest represents the request body for subscription
type SubscribeRequest struct {
	UserID   int64    `json:"user_id"`
	ItemID   int64    `json:"item_id"`
	Channels []string `json:"channels"`
}

// Handle handles POST /subscribe
func (h *SubscribeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.UserID <= 0 || req.ItemID <= 0 {
		http.Error(w, "user_id and item_id must be positive integers", http.StatusBadRequest)
		return
	}

	if len(req.Channels) == 0 {
		http.Error(w, "channels cannot be empty", http.StatusBadRequest)
		return
	}

	// Create subscription
	err := h.notificationService.Subscribe(req.UserID, req.ItemID, req.Channels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "subscribed"})
}
