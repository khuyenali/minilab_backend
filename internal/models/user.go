package models

import (
	"database/sql"
	"time"
	"mini-lab-api/internal/db"
)

// User represents a user in the system (domain model)
type User struct {
	ID        int32     `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name  string `json:"name" example:"John Doe Updated"`
	Email string `json:"email" example:"john.updated@example.com"`
}

// UserFilters represents filters for user queries
type UserFilters struct {
	Limit  int32  `json:"limit,omitempty"`
	Offset int32  `json:"offset,omitempty"`
	Search string `json:"search,omitempty"`
}

// ToDBUser converts domain model to database model
func (u *User) ToDBUser() db.User {
	return db.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		CreatedAt: sql.NullTime{Time: u.CreatedAt, Valid: !u.CreatedAt.IsZero()},
		UpdatedAt: sql.NullTime{Time: u.UpdatedAt, Valid: !u.UpdatedAt.IsZero()},
	}
}

// FromDBUser converts database model to domain model
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
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromDBUsers converts slice of database models to domain models
func FromDBUsers(dbUsers []db.User) []*User {
	users := make([]*User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = FromDBUser(dbUser)
	}
	return users
} 