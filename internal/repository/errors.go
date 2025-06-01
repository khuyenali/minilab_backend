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
	ErrInvalidUserIDs    = errors.New("one or more user IDs are invalid")
	ErrTaskNotFound     = errors.New("task not found")
	ErrSubTaskNotFound  = errors.New("sub-task not found")
	ErrAssignmentNotFound = errors.New("assignment not found")
) 