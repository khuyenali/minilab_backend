package models

import (
	"time"
	"mini-lab-api/internal/db"
	"database/sql"
)

// Task represents a task in the system (domain model)
type Task struct {
	ID        int32      `json:"id" example:"1"`
	TaskName  string     `json:"task_name" example:"Organize Company Retreat"`
	Status    string     `json:"status" example:"draft"`
	Priority  int32      `json:"priority" example:"1"`
	StartTime *time.Time `json:"start_time,omitempty" example:"2024-07-01T00:00:00Z"`
	EndTime   *time.Time `json:"end_time,omitempty" example:"2024-07-05T00:00:00Z"`
	Note      *string    `json:"note,omitempty" example:"Annual company-wide retreat"`
	SubTasks  []*SubTask `json:"sub_tasks,omitempty"`
	CreatedAt time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SubTaskRequest represents a sub-task in the create task request
type SubTaskRequest struct {
	TypeID  int32   `json:"type_id" binding:"required" example:"10"`
	UserIDs []int32 `json:"user_ids,omitempty" swaggertype:"array,integer" example:"3,7"`
}

// CreateTaskRequest represents the request to create a new task with sub-tasks
type CreateTaskRequest struct {
	TaskName  string             `json:"task_name" binding:"required" example:"Organize Company Retreat"`
	Status    string             `json:"status,omitempty" example:"draft"`
	Priority  int32              `json:"priority,omitempty" example:"1"`
	StartTime *time.Time         `json:"start_time,omitempty" example:"2024-07-01T00:00:00Z"`
	EndTime   *time.Time         `json:"end_time,omitempty" example:"2024-07-05T00:00:00Z"`
	Note      *string            `json:"note,omitempty" example:"Annual company-wide retreat"`
	SubTasks  []SubTaskRequest   `json:"sub_tasks,omitempty"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	TaskName  string     `json:"task_name" example:"Updated Task Name"`
	Status    string     `json:"status" example:"processing"`
	Priority  int32      `json:"priority" example:"1"`
	StartTime *time.Time `json:"start_time,omitempty" example:"2024-07-01T00:00:00Z"`
	EndTime   *time.Time `json:"end_time,omitempty" example:"2024-07-05T00:00:00Z"`
	Note      *string    `json:"note,omitempty" example:"Updated note"`
}

// FromDBTask converts database model to domain model
func FromDBTask(dbTask db.Task) *Task {
	var createdAt, updatedAt time.Time
	var startTime, endTime *time.Time
	var note *string
	
	if dbTask.CreatedAt.Valid {
		createdAt = dbTask.CreatedAt.Time
	}
	if dbTask.UpdatedAt.Valid {
		updatedAt = dbTask.UpdatedAt.Time
	}
	if dbTask.StartTime.Valid {
		startTime = &dbTask.StartTime.Time
	}
	if dbTask.EndTime.Valid {
		endTime = &dbTask.EndTime.Time
	}
	if dbTask.Note.Valid {
		note = &dbTask.Note.String
	}
	
	return &Task{
		ID:        dbTask.ID,
		TaskName:  dbTask.TaskName,
		Status:    string(dbTask.Status),
		Priority:  dbTask.Priority,
		StartTime: startTime,
		EndTime:   endTime,
		Note:      note,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromDBTasks converts slice of database models to domain models
func FromDBTasks(dbTasks []db.Task) []*Task {
	tasks := make([]*Task, len(dbTasks))
	for i, dbTask := range dbTasks {
		tasks[i] = FromDBTask(dbTask)
	}
	return tasks
}

// ToCreateTaskParams converts CreateTaskRequest to database parameters
func (req *CreateTaskRequest) ToCreateTaskParams() db.CreateTaskParams {
	// Set defaults
	status := req.Status
	if status == "" {
		status = "draft"
	}
	priority := req.Priority
	if priority == 0 {
		priority = 2 // default to medium
	}

	var startTime, endTime sql.NullTime
	var note sql.NullString

	if req.StartTime != nil {
		startTime = sql.NullTime{Time: *req.StartTime, Valid: true}
	}
	if req.EndTime != nil {
		endTime = sql.NullTime{Time: *req.EndTime, Valid: true}
	}
	if req.Note != nil {
		note = sql.NullString{String: *req.Note, Valid: true}
	}

	return db.CreateTaskParams{
		TaskName:  req.TaskName,
		Status:    db.TaskStatus(status),
		Priority:  priority,
		StartTime: startTime,
		EndTime:   endTime,
		Note:      note,
	}
} 