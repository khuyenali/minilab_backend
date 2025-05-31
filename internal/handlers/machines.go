package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mini-lab-api/internal/models"
	"mini-lab-api/internal/service"
)

// GetMachines handles GET /machine
// @Summary Get all machines
// @Description Get a list of all machines
// @Tags machines
// @Produce json
// @Success 200 {object} ApiResponse{data=[]models.Machine} "List of machines"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/machine [get]
func (h *Handler) GetMachines(c *gin.Context) {
	machines, err := h.machineService.GetMachines(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve machines",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    machines,
		"status":  "success",
		"message": "Machines retrieved successfully",
	})
}

// GetMachine handles GET /machine/:id
// @Summary Get machine by ID
// @Description Get a single machine by its ID
// @Tags machines
// @Produce json
// @Param id path int true "Machine ID"
// @Success 200 {object} ApiResponse{data=models.Machine} "Machine details"
// @Failure 400 {object} ApiResponse "Invalid machine ID"
// @Failure 404 {object} ApiResponse "Machine not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/machine/{id} [get]
func (h *Handler) GetMachine(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid machine ID",
			"status":  "error",
			"message": "Machine ID must be a number",
		})
		return
	}

	machine, err := h.machineService.GetMachine(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrMachineNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Machine not found",
				"status":  "error",
				"message": "Machine with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidMachineID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid machine ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve machine",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    machine,
		"status":  "success",
		"message": "Machine retrieved successfully",
	})
}

// CreateMachine handles POST /machine
// @Summary Create a new machine
// @Description Create a new machine with optional task type assignment
// @Tags machines
// @Accept json
// @Produce json
// @Param machine body models.CreateMachineRequest true "Machine creation data"
// @Success 201 {object} ApiResponse{data=models.Machine} "Machine created successfully"
// @Failure 400 {object} ApiResponse "Invalid input data"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/machine [post]
func (h *Handler) CreateMachine(c *gin.Context) {
	var req models.CreateMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	machine, err := h.machineService.CreateMachine(c.Request.Context(), req)
	if err != nil {
		// Handle validation errors
		if err == service.ErrInvalidMachineName || 
		   err == service.ErrMachineNameTooLong || 
		   err == service.ErrInvalidMachineQuantity ||
		   err == service.ErrInvalidMachineEstimateTime {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create machine",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    machine,
		"status":  "success",
		"message": "Machine created successfully",
	})
}

// UpdateMachine handles PUT /machine/:id
// @Summary Update a machine
// @Description Update an existing machine by ID
// @Tags machines
// @Accept json
// @Produce json
// @Param id path int true "Machine ID"
// @Param machine body models.UpdateMachineRequest true "Machine update data"
// @Success 200 {object} ApiResponse{data=models.Machine} "Machine updated successfully"
// @Failure 400 {object} ApiResponse "Invalid input data"
// @Failure 404 {object} ApiResponse "Machine not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/machine/{id} [put]
func (h *Handler) UpdateMachine(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid machine ID",
			"status":  "error",
			"message": "Machine ID must be a number",
		})
		return
	}

	var req models.UpdateMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	machine, err := h.machineService.UpdateMachine(c.Request.Context(), int32(id), req)
	if err != nil {
		if err == service.ErrMachineNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Machine not found",
				"status":  "error",
				"message": "Machine with the specified ID does not exist",
			})
			return
		}

		// Handle validation errors
		if err == service.ErrInvalidMachineID || 
		   err == service.ErrInvalidMachineName || 
		   err == service.ErrMachineNameTooLong || 
		   err == service.ErrInvalidMachineQuantity ||
		   err == service.ErrInvalidMachineEstimateTime {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation error",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update machine",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    machine,
		"status":  "success",
		"message": "Machine updated successfully",
	})
}

// DeleteMachine handles DELETE /machine/:id
// @Summary Delete a machine
// @Description Delete a machine by ID
// @Tags machines
// @Produce json
// @Param id path int true "Machine ID"
// @Success 200 {object} ApiResponse "Machine deleted successfully"
// @Failure 400 {object} ApiResponse "Invalid machine ID"
// @Failure 404 {object} ApiResponse "Machine not found"
// @Failure 500 {object} ApiResponse "Internal server error"
// @Router /api/v1/machine/{id} [delete]
func (h *Handler) DeleteMachine(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid machine ID",
			"status":  "error",
			"message": "Machine ID must be a number",
		})
		return
	}

	err = h.machineService.DeleteMachine(c.Request.Context(), int32(id))
	if err != nil {
		if err == service.ErrMachineNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Machine not found",
				"status":  "error",
				"message": "Machine with the specified ID does not exist",
			})
			return
		}

		if err == service.ErrInvalidMachineID {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid machine ID",
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete machine",
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Machine deleted successfully",
	})
} 