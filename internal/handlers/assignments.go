package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/service"
)

// CreateAssignment handles POST /assignments
// @Summary Create a new assignment
// @Description Create a new assignment for a user to a sub-task
// @Tags assignments
// @Accept json
// @Produce json
// @Param assignment body models.CreateAssignmentRequest true "Assignment creation data"
// @Success 201 {object} ApiResponse{data=models.UserSubTaskAssignment} "Assignment created successfully"
// @Failure 400 {object} ApiResponse "Invalid input data"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/assignments [post]
func (h *Handler) CreateAssignment(c *gin.Context) {
	var req models.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	assignment, err := h.assignmentService.CreateAssignment(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidSubTaskForAssignment || 
		   err == service.ErrInvalidUserForAssignment {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input data",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		if err == service.ErrAssignmentAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Assignment already exists",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create assignment",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    assignment,
		"status":  "success",
		"message": "Assignment created successfully",
	})
}

// DeleteAssignment handles DELETE /assignments/:id
// @Summary Delete an assignment
// @Description Delete an assignment by ID
// @Tags assignments
// @Produce json
// @Param id path int true "Assignment ID"
// @Success 200 {object} ApiResponse "Assignment deleted successfully"
// @Failure 400 {object} ApiResponse "Invalid assignment ID"
// @Failure 404 {object} ApiResponse "Assignment not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/assignments/{id} [delete]
func (h *Handler) DeleteAssignment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid assignment ID",
			"status":  "error",
			"message": "Assignment ID must be a number",
		})
		return
	}

	err = h.assignmentService.DeleteAssignment(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrAssignmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Assignment not found",
				"status":  "error",
				"message": "Assignment with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidAssignmentID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid assignment ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete assignment",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Assignment deleted successfully",
	})
} 