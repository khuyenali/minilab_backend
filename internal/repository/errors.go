package repository

import "errors"

// Repository errors
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserEmailExists  = errors.New("user with this email already exists")
	ErrRoleNotFound     = errors.New("role not found")
	ErrTaskTypeNotFound = errors.New("task type not found")
	ErrMachineNotFound  = errors.New("machine not found")
	ErrInvalidMachineIDs = errors.New("one or more machine IDs are invalid")
) 