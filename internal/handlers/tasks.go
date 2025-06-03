package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/service"
)

// GetTasks handles GET /tasks
// @Summary Get all tasks
// @Description Get a list of all tasks with sub-tasks containing their assignments
// @Tags tasks
// @Produce json
// @Success 200 {object} ApiResponse{data=[]models.Task} "List of tasks with sub-tasks and their assignments"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks [get]
func (h *Handler) GetTasks(c *gin.Context) {
	tasks, err := h.taskService.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve tasks",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    tasks,
		"status":  "success",
		"message": "Tasks retrieved successfully",
	})
}

// GetAvailableTasks handles GET /tasks/available
// @Summary Get available tasks for auto assignment
// @Description Get draft task IDs available for auto assignment based on priority, user loading constraints (8 hours max), and machine availability. Only tasks with status="draft" are considered. Tasks with pending/process status use machines, finished tasks release machines.
// @Tags tasks
// @Produce json
// @Success 200 {object} ApiResponse{data=[]int32} "Available task IDs for auto assignment"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks/available [get]
func (h *Handler) GetAvailableTasks(c *gin.Context) {
	taskIDs, err := h.taskService.GetAvailableTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get available tasks",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    taskIDs,
		"status":  "success",
		"message": "Available task IDs retrieved successfully",
	})
}

// GetTask handles GET /tasks/:id
// @Summary Get task by ID
// @Description Get a single task by its ID with sub-tasks containing their assignments
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} ApiResponse{data=models.Task} "Task details with sub-tasks and their assignments"
// @Failure 400 {object} ApiResponse "Invalid task ID"
// @Failure 404 {object} ApiResponse "Task not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid task ID",
			"status":  "error",
			"message": "Task ID must be a number",
		})
		return
	}

	task, err := h.taskService.GetTask(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task not found",
				"status":  "error",
				"message": "Task with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidTaskID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid task ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve task",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"status":  "success",
		"message": "Task retrieved successfully",
	})
}

// CreateTask handles POST /tasks
// @Summary Create a new task with sub-tasks and assignments
// @Description Create a new task with optional sub-tasks and user assignments
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Task creation data with sub-tasks and assignments"
// @Success 201 {object} ApiResponse{data=models.Task} "Task created successfully with sub-tasks and assignments"
// @Failure 400 {object} ApiResponse "Invalid input data"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks [post]
func (h *Handler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	task, err := h.taskService.CreateTask(c.Request.Context(), req)
	if err != nil {
		// Handle validation errors
		if err == service.ErrInvalidTaskName || 
		   err == service.ErrTaskNameTooLong || 
		   err == service.ErrTaskNoteTooLong ||
		   err == service.ErrInvalidTaskStatus ||
		   err == service.ErrInvalidTaskPriority {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		// Handle sub-task validation errors (these will have detailed messages)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    task,
		"status":  "success",
		"message": "Task created successfully with sub-tasks and assignments",
	})
}

// UpdateTask handles PUT /tasks/:id
// @Summary Update a task
// @Description Update an existing task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.UpdateTaskRequest true "Task update data"
// @Success 200 {object} ApiResponse{data=models.Task} "Task updated successfully"
// @Failure 400 {object} ApiResponse "Invalid input data"
// @Failure 404 {object} ApiResponse "Task not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks/{id} [put]
func (h *Handler) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid task ID",
			"status":  "error",
			"message": "Task ID must be a number",
		})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	task, err := h.taskService.UpdateTask(c.Request.Context(), int32(id), req)
	if err != nil {
		if err == service.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task not found",
				"status":  "error",
				"message": "Task with the specified ID does not exist",
			})
			return
		}

		// Handle validation errors
		if err == service.ErrInvalidTaskID || 
		   err == service.ErrInvalidTaskName || 
		   err == service.ErrTaskNameTooLong || 
		   err == service.ErrTaskNoteTooLong ||
		   err == service.ErrInvalidTaskStatus ||
		   err == service.ErrInvalidTaskPriority {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update task",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"status":  "success",
		"message": "Task updated successfully",
	})
}

// UpdateTaskStatusToPending handles PUT /tasks/:id/status
// @Summary Update task status from draft to pending
// @Description Update a task's status from draft to pending
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} ApiResponse{data=models.Task} "Task status updated successfully"
// @Failure 400 {object} ApiResponse "Invalid task ID or status transition"
// @Failure 404 {object} ApiResponse "Task not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks/{id}/status [put]
func (h *Handler) UpdateTaskStatusToPending(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid task ID",
			"status":  "error",
			"message": "Task ID must be a number",
		})
		return
	}

	task, err := h.taskService.UpdateTaskStatusToPending(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task not found",
				"status":  "error",
				"message": "Task with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidTaskID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid task ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		if err == service.ErrInvalidTaskStatusTransition {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid status transition",
				"status":  "error",
				"message": "Task must be in draft status to be updated to pending",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update task status",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"status":  "success",
		"message": "Task status updated to pending successfully",
	})
}

// DeleteTask handles DELETE /tasks/:id
// @Summary Delete a task
// @Description Delete an existing task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} ApiResponse "Task deleted successfully"
// @Failure 400 {object} ApiResponse "Invalid task ID"
// @Failure 404 {object} ApiResponse "Task not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid task ID",
			"status":  "error",
			"message": "Task ID must be a number",
		})
		return
	}

	err = h.taskService.DeleteTask(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task not found",
				"status":  "error",
				"message": "Task with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidTaskID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid task ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete task",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Task deleted successfully",
	})
} 