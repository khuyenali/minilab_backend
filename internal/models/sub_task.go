package models

import (
	"time"
	"mini-lab-api/internal/db"
	"database/sql"
)

// SubTask represents a sub-task in the system (domain model)
type SubTask struct {
	ID             int32                    `json:"id" example:"105"`
	TaskID         int32                    `json:"task_id" example:"3"`
	TypeID         int32                    `json:"type_id" example:"10"`
	SubTaskName    *string                  `json:"sub_task_name,omitempty" example:"Venue Booking Sub-Task"`
	Description    *string                  `json:"description,omitempty" example:"Book the venue for the retreat"`
	EstimateEffort *int32                   `json:"estimate_effort,omitempty" example:"8"`
	Assignments    []*UserSubTaskAssignment `json:"assignments,omitempty"`
}

// UserSubTaskAssignment represents the assignment of a user to a sub-task
type UserSubTaskAssignment struct {
	AssignmentID int32      `json:"assignment_id" example:"201"`
	UserID       int32      `json:"user_id" example:"3"`
	SubTaskID    int32      `json:"sub_task_id" example:"105"`
	Report       *string    `json:"report,omitempty" example:"Task assigned to implement the UI"`
	AssignedAt   time.Time  `json:"assigned_at" example:"2023-01-01T12:00:00Z"`
	User         *UserBasic `json:"user_details,omitempty"`
}

// CreateSubTaskRequest represents the request to create a new sub-task
type CreateSubTaskRequest struct {
	TaskID         int32   `json:"task_id" binding:"required" example:"1"`
	TypeID         int32   `json:"type_id" binding:"required" example:"10"`
	SubTaskName    *string `json:"sub_task_name,omitempty" example:"Design UI"`
	Description    *string `json:"description,omitempty" example:"Design the user interface"`
	EstimateEffort *int32  `json:"estimate_effort,omitempty" example:"8"`
}

// UpdateSubTaskRequest represents the request to update a sub-task
type UpdateSubTaskRequest struct {
	SubTaskName    *string `json:"sub_task_name,omitempty" example:"Updated Sub-Task Name"`
	Description    *string `json:"description,omitempty" example:"Updated description"`
	EstimateEffort *int32  `json:"estimate_effort,omitempty" example:"12"`
}

// CreateUserSubTaskAssignmentRequest represents the request to assign a user to a sub-task
type CreateUserSubTaskAssignmentRequest struct {
	UserID int32   `json:"user_id" binding:"required" example:"3"`
	Report *string `json:"report,omitempty" example:"User assigned to implement the feature"`
}

// UpdateUserSubTaskAssignmentRequest represents the request to update an assignment
type UpdateUserSubTaskAssignmentRequest struct {
	Report *string `json:"report,omitempty" example:"Progress update: 50% complete"`
}

// FromDBSubTask converts database model to domain model
func FromDBSubTask(dbSubTask db.SubTask) *SubTask {
	var subTaskName, description *string
	var estimateEffort *int32
	
	if dbSubTask.SubTaskName.Valid {
		subTaskName = &dbSubTask.SubTaskName.String
	}
	if dbSubTask.Description.Valid {
		description = &dbSubTask.Description.String
	}
	if dbSubTask.EstimateEffort.Valid {
		estimateEffort = &dbSubTask.EstimateEffort.Int32
	}
	
	return &SubTask{
		ID:             dbSubTask.ID,
		TaskID:         dbSubTask.TaskID,
		TypeID:         dbSubTask.TypeID,
		SubTaskName:    subTaskName,
		Description:    description,
		EstimateEffort: estimateEffort,
	}
}

// FromDBSubTasks converts slice of database models to domain models
func FromDBSubTasks(dbSubTasks []db.SubTask) []*SubTask {
	subTasks := make([]*SubTask, len(dbSubTasks))
	for i, dbSubTask := range dbSubTasks {
		subTasks[i] = FromDBSubTask(dbSubTask)
	}
	return subTasks
}

// FromDBUserSubTaskAssignment converts database model to domain model
func FromDBUserSubTaskAssignment(dbAssignment db.UserToSubTask) *UserSubTaskAssignment {
	var assignedAt time.Time
	var report *string
	
	if dbAssignment.AssignedAt.Valid {
		assignedAt = dbAssignment.AssignedAt.Time
	}
	if dbAssignment.Report.Valid {
		report = &dbAssignment.Report.String
	}
	
	return &UserSubTaskAssignment{
		AssignmentID: dbAssignment.ID,
		UserID:       dbAssignment.UserID,
		SubTaskID:    dbAssignment.SubTaskID,
		Report:       report,
		AssignedAt:   assignedAt,
	}
}

// FromDBUserSubTaskAssignments converts slice of database models to domain models
func FromDBUserSubTaskAssignments(dbAssignments []db.UserToSubTask) []*UserSubTaskAssignment {
	assignments := make([]*UserSubTaskAssignment, len(dbAssignments))
	for i, dbAssignment := range dbAssignments {
		assignments[i] = FromDBUserSubTaskAssignment(dbAssignment)
	}
	return assignments
}

// ToCreateSubTaskParams converts CreateSubTaskRequest to database parameters
func (req *CreateSubTaskRequest) ToCreateSubTaskParams() db.CreateSubTaskParams {
	var subTaskName, description sql.NullString
	var estimateEffort sql.NullInt32

	if req.SubTaskName != nil {
		subTaskName = sql.NullString{String: *req.SubTaskName, Valid: true}
	}
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.EstimateEffort != nil {
		estimateEffort = sql.NullInt32{Int32: *req.EstimateEffort, Valid: true}
	}

	return db.CreateSubTaskParams{
		TaskID:         req.TaskID,
		TypeID:         req.TypeID,
		SubTaskName:    subTaskName,
		Description:    description,
		EstimateEffort: estimateEffort,
	}
} 