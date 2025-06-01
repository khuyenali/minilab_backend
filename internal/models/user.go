package models

import (
	"database/sql"
	"time"
	"mini-lab-api/internal/db"
)

// User represents a user in the system (domain model)
type User struct {
	ID        int32            `json:"id" example:"1"`
	Name      string           `json:"name" example:"John Doe"`
	Email     string           `json:"email" example:"john@example.com"`
	RoleID    int32            `json:"-"` // Internal use only, not exposed in API responses
	Role      string           `json:"role" example:"member"`
	TaskTypes []*TaskTypeBasic `json:"task_types"`
	CreatedAt time.Time        `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time        `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// Role represents a role in the system (domain model)
type Role struct {
	ID        int32     `json:"id" example:"1"`
	RoleName  string    `json:"role_name" example:"admin"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Name   string `json:"name" binding:"required" example:"John Doe"`
	Email  string `json:"email" binding:"required,email" example:"john@example.com"`
	RoleID *int32 `json:"role_id" example:"3"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name   string `json:"name" example:"John Doe Updated"`
	Email  string `json:"email" example:"john.updated@example.com"`
	RoleID *int32 `json:"role_id" example:"2"`
}

// UserTaskAssignmentRequest represents the request to assign task types to a user
type UserTaskAssignmentRequest struct {
	TaskTypeIDs []int32 `json:"task_types" binding:"required" swaggertype:"array,integer" example:"1,2,3"`
}

// UserFilters represents filters for user queries
type UserFilters struct {
	Limit  int32  `json:"limit,omitempty"`
	Offset int32  `json:"offset,omitempty"`
	Search string `json:"search,omitempty"`
	RoleID *int32 `json:"role_id,omitempty"`
}

// ToDBUser converts domain model to database model
func (u *User) ToDBUser() db.User {
	return db.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		RoleID:    u.RoleID,
		CreatedAt: sql.NullTime{Time: u.CreatedAt, Valid: !u.CreatedAt.IsZero()},
		UpdatedAt: sql.NullTime{Time: u.UpdatedAt, Valid: !u.UpdatedAt.IsZero()},
	}
}

// FromDBUser converts database model to domain model (simple user without role name)
func FromDBUser(dbUser db.User) *User {
	var createdAt, updatedAt time.Time
	
	if dbUser.CreatedAt.Valid {
		createdAt = dbUser.CreatedAt.Time
	}
	if dbUser.UpdatedAt.Valid {
		updatedAt = dbUser.UpdatedAt.Time
	}
	
	return &User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		RoleID:    dbUser.RoleID,
		Role:      "", // Will be empty for simple conversions
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromGetUserRow converts query result to domain model (with role name)
func FromGetUserRow(row db.GetUserRow) *User {
	var createdAt, updatedAt time.Time
	
	if row.CreatedAt.Valid {
		createdAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time
	}
	
	return &User{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		RoleID:    row.RoleID,
		Role:      row.RoleName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromListUsersRow converts list query result to domain model
func FromListUsersRow(row db.ListUsersRow) *User {
	var createdAt, updatedAt time.Time
	
	if row.CreatedAt.Valid {
		createdAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time
	}
	
	return &User{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		RoleID:    row.RoleID,
		Role:      row.RoleName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromGetUserByEmailRow converts email query result to domain model
func FromGetUserByEmailRow(row db.GetUserByEmailRow) *User {
	var createdAt, updatedAt time.Time
	
	if row.CreatedAt.Valid {
		createdAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time
	}
	
	return &User{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		RoleID:    row.RoleID,
		Role:      row.RoleName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromGetUsersByRoleRows converts role query results to domain models
func FromGetUsersByRoleRows(rows []db.GetUsersByRoleRow) []*User {
	users := make([]*User, len(rows))
	for i, row := range rows {
		var createdAt, updatedAt time.Time
		
		if row.CreatedAt.Valid {
			createdAt = row.CreatedAt.Time
		}
		if row.UpdatedAt.Valid {
			updatedAt = row.UpdatedAt.Time
		}
		
		users[i] = &User{
			ID:        row.ID,
			Name:      row.Name,
			Email:     row.Email,
			RoleID:    row.RoleID,
			Role:      row.RoleName,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}
	return users
}

// FromGetUsersByRoleIDRows converts role ID query results to domain models
func FromGetUsersByRoleIDRows(rows []db.GetUsersByRoleIDRow) []*User {
	users := make([]*User, len(rows))
	for i, row := range rows {
		var createdAt, updatedAt time.Time
		
		if row.CreatedAt.Valid {
			createdAt = row.CreatedAt.Time
		}
		if row.UpdatedAt.Valid {
			updatedAt = row.UpdatedAt.Time
		}
		
		users[i] = &User{
			ID:        row.ID,
			Name:      row.Name,
			Email:     row.Email,
			RoleID:    row.RoleID,
			Role:      row.RoleName,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}
	return users
}

// FromListUsersRows converts list query results to domain models
func FromListUsersRows(rows []db.ListUsersRow) []*User {
	users := make([]*User, len(rows))
	for i, row := range rows {
		users[i] = FromListUsersRow(row)
	}
	return users
}

// FromDBRole converts database role model to domain model
func FromDBRole(dbRole db.Role) *Role {
	var createdAt, updatedAt time.Time
	
	if dbRole.CreatedAt.Valid {
		createdAt = dbRole.CreatedAt.Time
	}
	if dbRole.UpdatedAt.Valid {
		updatedAt = dbRole.UpdatedAt.Time
	}
	
	return &Role{
		ID:        dbRole.ID,
		RoleName:  dbRole.RoleName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromDBRoles converts slice of database role models to domain models
func FromDBRoles(dbRoles []db.Role) []*Role {
	roles := make([]*Role, len(dbRoles))
	for i, dbRole := range dbRoles {
		roles[i] = FromDBRole(dbRole)
	}
	return roles
}

// FromGetUserTaskTypesRows converts user task types query results to domain models
func FromGetUserTaskTypesRows(rows []db.GetUserTaskTypesRow) []*TaskTypeBasic {
	taskTypes := make([]*TaskTypeBasic, len(rows))
	for i, row := range rows {
		var description *string
		var createdAt, updatedAt time.Time
		
		if row.Description.Valid {
			description = &row.Description.String
		}
		if row.CreatedAt.Valid {
			createdAt = row.CreatedAt.Time
		}
		if row.UpdatedAt.Valid {
			updatedAt = row.UpdatedAt.Time
		}
		
		taskTypes[i] = &TaskTypeBasic{
			ID:          row.TypeID,
			TypeName:    row.TypeName,
			Description: description,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
	}
	return taskTypes
} 