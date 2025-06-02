package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/service"
)

// UpdateAssignmentToProcessing handles PUT /assignments/process/:id
// @Summary Change assignment status from pending to processing
// @Description Change assignment status from pending to processing
// @Tags assignments
// @Param id path int true "Assignment ID"
// @Success 200 {object} ApiResponse{data=models.UserSubTaskAssignment} "Assignment status updated to processing successfully"
// @Failure 400 {object} ApiResponse "Invalid assignment ID or invalid status transition"
// @Failure 404 {object} ApiResponse "Assignment not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/assignments/process/{id} [put]
func (h *Handler) UpdateAssignmentToProcessing(c *gin.Context) {
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

	assignment, err := h.assignmentService.UpdateAssignmentToProcessing(c.Request.Context(), int32(id))
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

		if err == service.ErrAssignmentNotInPendingStatus {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid status transition",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update assignment status",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    assignment,
		"status":  "success",
		"message": "Assignment status updated to processing successfully",
	})
}

// FinishAssignment handles PUT /assignments/finish/:id
// @Summary Change assignment status from processing to finish with report
// @Description Change assignment status from processing to finish with a required report
// @Tags assignments
// @Accept json
// @Param id path int true "Assignment ID"
// @Param assignment body models.FinishAssignmentRequest true "Assignment finish data with report"
// @Success 200 {object} ApiResponse{data=models.UserSubTaskAssignment} "Assignment status updated to finish successfully"
// @Failure 400 {object} ApiResponse "Invalid assignment ID, invalid status transition, or missing report"
// @Failure 404 {object} ApiResponse "Assignment not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/assignments/finish/{id} [put]
func (h *Handler) FinishAssignment(c *gin.Context) {
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

	var req models.FinishAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	assignment, err := h.assignmentService.FinishAssignment(c.Request.Context(), int32(id), req.Report)
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

		if err == service.ErrAssignmentNotInProcessingStatus {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid status transition",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		if err == service.ErrAssignmentReportRequired {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Report required",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to finish assignment",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    assignment,
		"status":  "success",
		"message": "Assignment finished successfully",
	})
} 