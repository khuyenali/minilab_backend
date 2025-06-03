// Package main Mini Lab API
//
// A simple RESTful API server built with Go, Gin, sqlc, and golang-migrate.
//
// Terms Of Service: http://swagger.io/terms/
//
// Schemes: http, https
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// License: MIT http://opensource.org/licenses/MIT
// Contact: Mini Lab API <support@minilab.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "github.com/prometheus/client_golang/prometheus"        // Added for collector registration
	// "github.com/prometheus/client_golang/prometheus/collectors" // Added for Go and Process collectors
	// "github.com/prometheus/client_golang/prometheus/promhttp" // No longer directly used
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/zsais/go-gin-prometheus" // Import for Gin Prometheus middleware
	"mini-lab-api/internal/config"
	"mini-lab-api/internal/handlers"
	"mini-lab-api/internal/repository"
	"mini-lab-api/internal/service"
	
	_ "mini-lab-api/docs" // Import generated docs
)

// @title Mini Lab API
// @version 1.0
// @description A simple RESTful API server built with Go, Gin, sqlc, and golang-migrate.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.minilab.com/support
// @contact.email support@minilab.com

// @license.name MIT
// @license.url http://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Register standard Go collectors and process collector
	// These will be exposed alongside ginprometheus metrics if it uses the default registry.
	// prometheus.MustRegister(collectors.NewGoCollector()) // Removed due to duplicate registration panic
	// prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})) // Removed due to duplicate registration panic

	// Initialize configuration
	cfg := config.New()

	// Initialize database connection
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Initialize repository layer
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	machineRepo := repository.NewMachineRepository(db)
	taskTypeRepo := repository.NewTaskTypeRepository(db, machineRepo, userRepo)
	userToTypeRepo := repository.NewUserToTypeRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Initialize service layer
	userService := service.NewUserService(userRepo, userToTypeRepo, taskTypeRepo)
	roleService := service.NewRoleService(roleRepo)
	taskTypeService := service.NewTaskTypeService(taskTypeRepo, userRepo)
	machineService := service.NewMachineService(machineRepo, taskTypeRepo)
	taskService := service.NewTaskService(taskRepo, taskTypeRepo, userRepo, db)
	assignmentService := service.NewAssignmentService(db, taskRepo)

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	r := gin.Default()

	// Initialize go-gin-prometheus middleware
	// This will collect metrics AND expose /metrics by default with this library version
	p := ginprometheus.NewPrometheus("gin") 
	p.Use(r) 

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Initialize handlers with dependencies
	h := handlers.NewHandler(userService, roleService, taskTypeService, machineService, taskService, assignmentService)

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", healthCheck)

	// API routes
	api := r.Group("/api/v1")
	{
		// Example endpoints - you can expand these
		api.GET("/ping", h.Ping)
		
		// Role routes
		roles := api.Group("/roles")
		{
			roles.GET("", h.GetRoles)
			roles.GET("/:id", h.GetRole)
		}
		
		// User routes
		users := api.Group("/users")
		{
			users.GET("", h.GetUsers)
			users.POST("", h.CreateUser)
			users.GET("/:id", h.GetUser)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}
		
		// User task assignment route (members only)
		api.PUT("/user/:id", h.AssignUserTaskTypes)
		
		// Task type routes
		taskTypes := api.Group("/task_type")
		{
			taskTypes.GET("", h.GetTaskTypes)
			taskTypes.POST("", h.CreateTaskType)
			taskTypes.GET("/:id", h.GetTaskType)
			taskTypes.PUT("/:id", h.UpdateTaskType)
			taskTypes.DELETE("/:id", h.DeleteTaskType)
		}
		
		// Machine routes
		machines := api.Group("/machine")
		{
			machines.GET("", h.GetMachines)
			machines.POST("", h.CreateMachine)
			machines.GET("/:id", h.GetMachine)
			machines.PUT("/:id", h.UpdateMachine)
			machines.DELETE("/:id", h.DeleteMachine)
		}
		
		// Task cleanup route (separate path to avoid conflicts)
		api.DELETE("/cleanup-draft-assignments", h.CleanDraftTaskAssignments)
		
		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.GET("", h.GetTasks)
			tasks.GET("/available", h.GetAvailableTasks)
			tasks.POST("", h.CreateTask)
			tasks.GET("/:id", h.GetTask)
			tasks.PUT("/:id", h.UpdateTask)
			tasks.PUT("/:id/pending", h.UpdateTaskToPending)
			tasks.PUT("/:id/processing", h.UpdateTaskToProcessing)
			tasks.PUT("/:id/finish", h.FinishTask)
			tasks.DELETE("/:id", h.DeleteTask)
		}
		
		// Assignment routes (simplified)
		assignments := api.Group("/assignments")
		{
			assignments.POST("", h.CreateAssignment)
			assignments.DELETE("/:id", h.DeleteAssignment)
		}
	}

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Mini Lab API server on port %s", port)
	log.Printf("Swagger documentation available at: http://localhost:%s/swagger/index.html", port)
	log.Fatal(r.Run(":" + port))
}

// healthCheck handles the health check endpoint
// @Summary Health Check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{} "API is running"
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Mini Lab API is running",
		"version": "1.0.0",
	})
}

// corsMiddleware handles CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} 