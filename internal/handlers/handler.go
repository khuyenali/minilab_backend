package handlers

import (
	"mini-lab-api/internal/service"
)

// Handler struct holds dependencies for handlers
type Handler struct {
	userService service.UserService
}

// New creates a new handler instance
func New(userService service.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

// ApiResponse represents a standard API response
type ApiResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Error   string      `json:"error,omitempty" example:"Error message"`
} 