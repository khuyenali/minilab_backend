package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-lab-api/internal/service"
)

// GetRoles handles GET /roles
// @Summary Get all roles
// @Description Get a list of all available roles
// @Tags roles
// @Produce json
// @Success 200 {object} ApiResponse{data=[]models.Role} "List of roles"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/roles [get]
func (h *Handler) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve roles",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    roles,
		"status":  "success",
		"message": "Roles retrieved successfully",
	})
}

// GetRole handles GET /roles/:id
// @Summary Get role by ID
// @Description Get a single role by its ID
// @Tags roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} ApiResponse{data=models.Role} "Role details"
// @Failure 400 {object} ApiResponse "Invalid role ID"
// @Failure 404 {object} ApiResponse "Role not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/roles/{id} [get]
func (h *Handler) GetRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid role ID",
			"status":  "error",
			"message": "Role ID must be a number",
		})
		return
	}

	role, err := h.roleService.GetRole(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrRoleNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Role not found",
				"status":  "error",
				"message": "Role with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidRoleID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid role ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve role",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    role,
		"status":  "success",
		"message": "Role retrieved successfully",
	})
} 