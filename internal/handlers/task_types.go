package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/service"
)

// GetTaskTypes handles GET /task_type
// @Summary Get all task types
// @Description Get a list of all task types
// @Tags task-types
// @Produce json
// @Success 200 {object} ApiResponse{data=[]models.TaskType} "List of task types"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/task_type [get]
func (h *Handler) GetTaskTypes(c *gin.Context) {
	taskTypes, err := h.taskTypeService.GetTaskTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve task types",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    taskTypes,
		"status":  "success",
		"message": "Task types retrieved successfully",
	})
}

// GetTaskType handles GET /task_type/:id
// @Summary Get task type by ID
// @Description Get a single task type by its ID with associated machines
// @Tags task-types
// @Produce json
// @Param id path int true "Task Type ID"
// @Success 200 {object} ApiResponse{data=models.TaskType} "Task type details with machines"
// @Failure 400 {object} ApiResponse "Invalid task type ID"
// @Failure 404 {object} ApiResponse "Task type not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/task_type/{id} [get]
func (h *Handler) GetTaskType(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid task type ID",
			"status":  "error",
			"message": "Task type ID must be a number",
		})
		return
	}

	taskType, err := h.taskTypeService.GetTaskType(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrTaskTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task type not found",
				"status":  "error",
				"message": "Task type with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidTaskTypeID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid task type ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve task type",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    taskType,
		"status":  "success",
		"message": "Task type retrieved successfully",
	})
}

// CreateTaskType handles POST /task_type
// @Summary Create a new task type
// @Description Create a new task type with optional machine assignments
// @Tags task-types
// @Accept json
// @Produce json
// @Param task_type body models.CreateTaskTypeRequest true "Task type creation data"
// @Success 201 {object} ApiResponse{data=models.TaskType} "Task type created successfully"
// @Failure 400 {object} ApiResponse "Invalid input data"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/task_type [post]
func (h *Handler) CreateTaskType(c *gin.Context) {
	var req models.CreateTaskTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	taskType, err := h.taskTypeService.CreateTaskType(c.Request.Context(), req)
	if err != nil {
		// Handle validation errors
		if err == service.ErrInvalidTaskTypeName || 
		   err == service.ErrTaskTypeNameTooLong || 
		   err == service.ErrTaskTypeDescriptionTooLong ||
		   err == service.ErrInvalidMachineIDs {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create task type",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    taskType,
		"status":  "success",
		"message": "Task type created successfully",
	})
}

// UpdateTaskType handles PUT /task_type/:id
// @Summary Update a task type
// @Description Update an existing task type by ID
// @Tags task-types
// @Accept json
// @Produce json
// @Param id path int true "Task Type ID"
// @Param task_type body models.UpdateTaskTypeRequest true "Task type update data"
// @Success 200 {object} ApiResponse{data=models.TaskType} "Task type updated successfully"
// @Failure 400 {object} ApiResponse "Invalid input data"
// @Failure 404 {object} ApiResponse "Task type not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/task_type/{id} [put]
func (h *Handler) UpdateTaskType(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid task type ID",
			"status":  "error",
			"message": "Task type ID must be a number",
		})
		return
	}

	var req models.UpdateTaskTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	taskType, err := h.taskTypeService.UpdateTaskType(c.Request.Context(), int32(id), req)
	if err != nil {
		if err == service.ErrTaskTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task type not found",
				"status":  "error",
				"message": "Task type with the specified ID does not exist",
			})
			return
		}

		// Handle validation errors
		if err == service.ErrInvalidTaskTypeID || 
		   err == service.ErrInvalidTaskTypeName || 
		   err == service.ErrTaskTypeNameTooLong || 
		   err == service.ErrTaskTypeDescriptionTooLong ||
		   err == service.ErrInvalidMachineIDs {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update task type",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    taskType,
		"status":  "success",
		"message": "Task type updated successfully",
	})
}

// DeleteTaskType handles DELETE /task_type/:id
// @Summary Delete a task type
// @Description Delete a task type by ID
// @Tags task-types
// @Produce json
// @Param id path int true "Task Type ID"
// @Success 200 {object} ApiResponse "Task type deleted successfully"
// @Failure 400 {object} ApiResponse "Invalid task type ID"
// @Failure 404 {object} ApiResponse "Task type not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/task_type/{id} [delete]
func (h *Handler) DeleteTaskType(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid task type ID",
			"status":  "error",
			"message": "Task type ID must be a number",
		})
		return
	}

	err = h.taskTypeService.DeleteTaskType(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrTaskTypeNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task type not found",
				"status":  "error",
				"message": "Task type with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidTaskTypeID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid task type ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete task type",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Task type deleted successfully",
	})
} 