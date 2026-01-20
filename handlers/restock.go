package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kavya/noti/service"
)

// RestockHandler handles restock events
type RestockHandler struct {
	notificationService *service.NotificationService
}

// NewRestockHandler creates a new restock handler
func NewRestockHandler(notificationService *service.NotificationService) *RestockHandler {
	return &RestockHandler{
		notificationService: notificationService,
	}
}

// RestockRequest represents the request body for restock
type RestockRequest struct {
	ItemID   int64 `json:"item_id"`
	NewStock int   `json:"new_stock"`
}

// Handle handles POST /inventory/restock
func (h *RestockHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RestockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.ItemID <= 0 {
		http.Error(w, "item_id must be a positive integer", http.StatusBadRequest)
		return
	}

	// Process restock - this sends notifications to all subscribed users
	// Idempotent: if called multiple times, only PENDING subscriptions are processed
	err := h.notificationService.ProcessRestock(req.ItemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "restocked", "notifications_sent": "true"})
}
