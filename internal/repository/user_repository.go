package repository

import (
	"context"
	"database/sql"
	"mini-lab-api/internal/db"
	"mini-lab-api/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByID(ctx context.Context, id int32) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context, filters models.UserFilters) ([]*models.User, error)
	Create(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
	Update(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id int32) error
}

// userRepository implements UserRepository using sqlc generated code
type userRepository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewUserRepository creates a new user repository
func NewUserRepository(database *sql.DB) UserRepository {
	return &userRepository{
		db:      database,
		queries: db.New(database),
	}
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int32) (*models.User, error) {
	dbUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return models.FromDBUser(dbUser), nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return models.FromDBUser(dbUser), nil
}

// List retrieves users with optional filtering
func (r *userRepository) List(ctx context.Context, filters models.UserFilters) ([]*models.User, error) {
	// For now, we'll just get all users since we don't have complex filtering in our SQL
	// In a real app, you'd add more sophisticated query methods
	dbUsers, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	
	return models.FromDBUsers(dbUsers), nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	// Check if user with email already exists
	_, err := r.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrUserEmailExists
	}
	if err != ErrUserNotFound {
		return nil, err
	}
	
	dbUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}
	
	return models.FromDBUser(dbUser), nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.User, error) {
	// Check if user exists
	existingUser, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Prepare update params with existing values as defaults
	updateParams := db.UpdateUserParams{
		ID:    id,
		Name:  existingUser.Name,
		Email: existingUser.Email,
	}
	
	// Update only provided fields
	if req.Name != "" {
		updateParams.Name = req.Name
	}
	if req.Email != "" {
		// Check if email is already taken by another user
		if req.Email != existingUser.Email {
			existingEmailUser, err := r.GetByEmail(ctx, req.Email)
			if err == nil && existingEmailUser.ID != id {
				return nil, ErrUserEmailExists
			}
			if err != ErrUserNotFound {
				return nil, err
			}
		}
		updateParams.Email = req.Email
	}
	
	dbUser, err := r.queries.UpdateUser(ctx, updateParams)
	if err != nil {
		return nil, err
	}
	
	return models.FromDBUser(dbUser), nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id int32) error {
	// Check if user exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	return r.queries.DeleteUser(ctx, id)
} 
