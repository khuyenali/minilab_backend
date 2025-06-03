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
// @Description Get draft task IDs available for auto assignment and automatically create assignments. Based on priority, machine availability, and user task type matching. Only draft tasks are considered. Tasks are filtered by: 1) enough available machines (each sub-task consumes 1 machine), 2) users not already working on pending/processing/draft tasks, 3) users with matching task types, 4) sub-tasks that don't already have assignments. Automatically creates assignments for available users and returns task IDs. Use DELETE /tasks/assignments/draft first to clean existing assignments for a fresh start.
// @Tags tasks
// @Produce json
// @Success 200 {object} ApiResponse{data=[]int32} "Available task IDs with auto-created assignments"
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

// FinishTask handles PUT /tasks/:id/finish
// @Summary Finish a task with a required report
// @Description Finish a task by updating status to 'finish' with a required report
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.FinishTaskRequest true "Task finish data with report"
// @Success 200 {object} ApiResponse{data=models.Task} "Task finished successfully"
// @Failure 400 {object} ApiResponse "Invalid task ID, missing report, or invalid transition"
// @Failure 404 {object} ApiResponse "Task not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks/{id}/finish [put]
func (h *Handler) FinishTask(c *gin.Context) {
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

	var req models.FinishTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	task, err := h.taskService.FinishTask(c.Request.Context(), int32(id), req)
	if err != nil {
		if err == service.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task not found",
				"status":  "error",
				"message": "Task with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidTaskID || 
		   err == service.ErrTaskReportRequired ||
		   err == service.ErrInvalidTaskStatusTransition {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to finish task",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"status":  "success",
		"message": "Task finished successfully",
	})
}

// UpdateTaskToPending handles PUT /tasks/:id/pending
// @Summary Update task status to pending
// @Description Update a task status from draft to pending
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} ApiResponse{data=models.Task} "Task status updated to pending successfully"
// @Failure 400 {object} ApiResponse "Invalid task ID or invalid status transition"
// @Failure 404 {object} ApiResponse "Task not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks/{id}/pending [put]
func (h *Handler) UpdateTaskToPending(c *gin.Context) {
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

		if err == service.ErrInvalidTaskID || 
		   err == service.ErrInvalidTaskStatusTransition {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
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

// UpdateTaskToProcessing handles PUT /tasks/:id/processing
// @Summary Update task status to processing
// @Description Update a task status from pending to processing
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} ApiResponse{data=models.Task} "Task status updated to processing successfully"
// @Failure 400 {object} ApiResponse "Invalid task ID or invalid status transition"
// @Failure 404 {object} ApiResponse "Task not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/tasks/{id}/processing [put]
func (h *Handler) UpdateTaskToProcessing(c *gin.Context) {
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

	// Create an update request for processing status
	req := models.UpdateTaskStatusRequest{
		Status: "processing",
	}

	task, err := h.taskService.UpdateTaskStatus(c.Request.Context(), int32(id), req)
	if err != nil {
		if err == service.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Task not found",
				"status":  "error",
				"message": "Task with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidTaskID || 
		   err == service.ErrInvalidTaskStatusTransition {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
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
		"message": "Task status updated to processing successfully",
	})
}

// CleanDraftTaskAssignments handles DELETE /cleanup-draft-assignments
// @Summary Clean all assignments from draft tasks
// @Description Delete all existing assignments for tasks in draft status to prepare for clean auto-assignment
// @Tags tasks
// @Success 200 {object} ApiResponse "Draft task assignments cleaned successfully"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/cleanup-draft-assignments [delete]
func (h *Handler) CleanDraftTaskAssignments(c *gin.Context) {
	err := h.taskService.CleanDraftTaskAssignments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to clean draft task assignments",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Draft task assignments cleaned successfully",
	})
} 