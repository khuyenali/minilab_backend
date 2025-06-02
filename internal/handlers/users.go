package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/service"
)

// GetUsers handles GET /users
// @Summary Get all users
// @Description Get a list of all users with optional role filtering
// @Tags users
// @Produce json
// @Param role_id query int false "Filter by role ID" example(3)
// @Param limit query int false "Limit the number of results" example(50)
// @Param offset query int false "Offset for pagination" example(0)
// @Success 200 {object} ApiResponse{data=[]models.User} "List of users"
// @Failure 400 {object} ApiResponse "Invalid query parameters"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	// Parse query parameters
	filters := models.UserFilters{}
	
	// Parse role_id filter
	if roleIDStr := c.Query("role_id"); roleIDStr != "" {
		if roleID, err := strconv.ParseInt(roleIDStr, 10, 32); err == nil && roleID > 0 {
			roleID32 := int32(roleID)
			filters.RoleID = &roleID32
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid role_id parameter",
				"status":  "error",
				"message": "Role ID must be a positive number",
			})
			return
		}
	}
	
	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 32); err == nil && limit > 0 {
			filters.Limit = int32(limit)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid limit parameter",
				"status":  "error",
				"message": "Limit must be a positive number",
			})
			return
		}
	}
	
	// Parse offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && offset >= 0 {
			filters.Offset = int32(offset)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid offset parameter",
				"status":  "error",
				"message": "Offset must be a non-negative number",
			})
			return
		}
	}
	
	// Parse search
	if search := c.Query("search"); search != "" {
		filters.Search = search
	}

	users, err := h.userService.GetUsers(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve users",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    users,
		"status":  "success",
		"message": "Users retrieved successfully",
	})
}

// GetUser handles GET /users/:id
// @Summary Get user by ID
// @Description Get a single user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} ApiResponse{data=models.User} "User details"
// @Failure 400 {object} ApiResponse "Invalid user ID"
// @Failure 404 {object} ApiResponse "User not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"status":  "error",
			"message": "User ID must be a number",
		})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"status":  "error",
				"message": "User with the specified ID does not exist",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve user",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"status":  "success",
		"message": "User retrieved successfully",
	})
}

// CreateUser handles POST /users
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User information"
// @Success 201 {object} ApiResponse{data=models.User} "User created successfully"
// @Failure 400 {object} ApiResponse "Invalid request body"
// @Failure 409 {object} ApiResponse "User email already exists"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrUserEmailExists {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Email already exists",
				"status":  "error",
				"message": "A user with this email already exists",
			})
			return
		}

		// Handle validation errors
		if err == service.ErrInvalidUserName || err == service.ErrInvalidUserEmail ||
		   err == service.ErrUserNameTooLong || err == service.ErrUserEmailTooLong {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		// Handle role validation error
		if err == service.ErrInvalidUserRole {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user role",
				"status":  "error",
				"message": "Only members can be assigned task types",
			})
			return
		}

		// Check for invalid task type IDs error
		if strings.Contains(err.Error(), "invalid task type IDs:") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid task type IDs",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    user,
		"status":  "success",
		"message": "User created successfully",
	})
}

// UpdateUser handles PUT /users/:id
// @Summary Update user
// @Description Update an existing user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "Updated user information"
// @Success 200 {object} ApiResponse{data=models.User} "User updated successfully"
// @Failure 400 {object} ApiResponse "Invalid user ID or request body"
// @Failure 404 {object} ApiResponse "User not found"
// @Failure 409 {object} ApiResponse "User email already exists"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"status":  "error",
			"message": "User ID must be a number",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), int32(id), req)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"status":  "error",
				"message": "User with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrUserEmailExists {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Email already exists",
				"status":  "error",
				"message": "A user with this email already exists",
			})
			return
		}

		// Handle validation errors
		if err == service.ErrInvalidUserID || err == service.ErrInvalidUserName || 
		   err == service.ErrInvalidUserEmail || err == service.ErrUserNameTooLong || 
		   err == service.ErrUserEmailTooLong {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"status":  "success",
		"message": "User updated successfully",
	})
}

// DeleteUser handles DELETE /users/:id
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} ApiResponse "User deleted successfully"
// @Failure 400 {object} ApiResponse "Invalid user ID"
// @Failure 404 {object} ApiResponse "User not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"status":  "error",
			"message": "User ID must be a number",
		})
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"status":  "error",
				"message": "User with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidUserID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully",
	})
}

// AssignUserTaskTypes handles PUT /user/:id
// @Summary Assign task types to a user
// @Description Assign task types to a user (only for members)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param task_types body models.UserTaskAssignmentRequest true "Task type assignment data"
// @Success 200 {object} ApiResponse "Task types assigned successfully"
// @Failure 400 {object} ApiResponse "Invalid input data or user role"
// @Failure 404 {object} ApiResponse "User not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/user/{id} [put]
func (h *Handler) AssignUserTaskTypes(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"status":  "error",
			"message": "User ID must be a number",
		})
		return
	}

	var req models.UserTaskAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = h.userService.AssignTaskTypes(c.Request.Context(), int32(id), req)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"status":  "error",
				"message": "User with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidUserRole {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user role",
				"status":  "error",
				"message": "Only members can be assigned task types",
			})
			return
		}

		if err == service.ErrInvalidUserID || err == service.ErrInvalidTaskTypeAssignment {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		// Check for detailed invalid task type IDs error
		if strings.Contains(err.Error(), "invalid task type IDs:") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid task type IDs",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to assign task types",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Task types assigned successfully",
	})
} 
