package models

import (
	"time"
	"mini-lab-api/internal/db"
)

// UserBasic represents a user without full details (for nested contexts)
type UserBasic struct {
	ID   int32  `json:"user_id" example:"1"`
	Name string `json:"user_name" example:"John Doe"`
}

// TaskType represents a task type in the system (domain model)
type TaskType struct {
	ID          int32            `json:"id" example:"1"`
	TypeName    string           `json:"name" example:"3D Printing"`
	Description *string          `json:"description" example:"3D printing services"`
	Machines    []*MachineBasic  `json:"machines"`
	Users       []*UserBasic     `json:"users"`
	CreatedAt   time.Time        `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time        `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// TaskTypeBasic represents a task type without machines (for user contexts)
type TaskTypeBasic struct {
	ID          int32     `json:"id" example:"1"`
	TypeName    string    `json:"name" example:"3D Printing"`
	Description *string   `json:"description" example:"3D printing services"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// TaskTypeMinimal represents a task type without machines, created_at, and updated_at (for user contexts)
type TaskTypeMinimal struct {
	ID          int32   `json:"id" example:"1"`
	TypeName    string  `json:"name" example:"3D Printing"`
	Description *string `json:"description" example:"3D printing services"`
}

// CreateTaskTypeRequest represents the request to create a new task type
type CreateTaskTypeRequest struct {
	Name        string  `json:"name" binding:"required" example:"3D Printing"`
	Description *string `json:"description,omitempty" example:"3D printing and modeling tasks"`
	MachineIDs  []int32 `json:"machines,omitempty" swaggertype:"array,integer" example:"1,2"`
	UserIDs     []int32 `json:"users,omitempty" swaggertype:"array,integer" example:"1,2,3"`
}

// UpdateTaskTypeRequest represents the request to update a task type
type UpdateTaskTypeRequest struct {
	Name        string  `json:"name" example:"3D Printing Updated"`
	Description *string `json:"description,omitempty" example:"Updated description"`
	MachineIDs  []int32 `json:"machines,omitempty" swaggertype:"array,integer" example:"1,2,3"`
	UserIDs     []int32 `json:"users,omitempty" swaggertype:"array,integer" example:"1,2,3"`
}

// FromDBTaskType converts database model to domain model
func FromDBTaskType(dbTaskType db.TaskType) *TaskType {
	var createdAt, updatedAt time.Time
	var description *string
	
	if dbTaskType.CreatedAt.Valid {
		createdAt = dbTaskType.CreatedAt.Time
	}
	if dbTaskType.UpdatedAt.Valid {
		updatedAt = dbTaskType.UpdatedAt.Time
	}
	if dbTaskType.Description.Valid {
		description = &dbTaskType.Description.String
	}
	
	return &TaskType{
		ID:          dbTaskType.ID,
		TypeName:    dbTaskType.TypeName,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// FromDBTaskTypes converts slice of database models to domain models
func FromDBTaskTypes(dbTaskTypes []db.TaskType) []*TaskType {
	taskTypes := make([]*TaskType, len(dbTaskTypes))
	for i, dbTaskType := range dbTaskTypes {
		taskTypes[i] = FromDBTaskType(dbTaskType)
	}
	return taskTypes
}

// FromListTaskTypesRows converts list query results to domain models
func FromListTaskTypesRows(rows []db.TaskType) []*TaskType {
	taskTypes := make([]*TaskType, len(rows))
	for i, row := range rows {
		taskTypes[i] = FromDBTaskType(row)
	}
	return taskTypes
}

// FromGetTaskTypeRow converts get query result to domain model
func FromGetTaskTypeRow(row db.TaskType) *TaskType {
	return FromDBTaskType(row)
}

// FromGetTaskTypeUsersRows converts task type users query results to UserBasic models
func FromGetTaskTypeUsersRows(rows []db.GetTaskTypeUsersRow) []*UserBasic {
	users := make([]*UserBasic, len(rows))
	for i, row := range rows {
		users[i] = &UserBasic{
			ID:   row.UserID,
			Name: row.Name,
		}
	}
	return users
} 