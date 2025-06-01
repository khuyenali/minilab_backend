package handlers

import (
	"mini-lab-api/internal/service"
)

// Handler holds all the dependencies for handlers
type Handler struct {
	userService     service.UserService
	roleService     service.RoleService
	taskTypeService service.TaskTypeService
	machineService  service.MachineService
	taskService     service.TaskService
}

// NewHandler creates a new handler instance
func NewHandler(userService service.UserService, roleService service.RoleService, taskTypeService service.TaskTypeService, machineService service.MachineService, taskService service.TaskService) *Handler {
	return &Handler{
		userService:     userService,
		roleService:     roleService,
		taskTypeService: taskTypeService,
		machineService:  machineService,
		taskService:     taskService,
	}
}

// ApiResponse represents a standard API response
type ApiResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Error   string      `json:"error,omitempty" example:"Error message"`
} 