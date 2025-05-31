package repository

import "errors"

// Repository errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserEmailExists = errors.New("user with this email already exists")
) 