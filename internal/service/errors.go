package service

import "errors"

// Service errors
var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserEmailExists   = errors.New("user with this email already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidUserName   = errors.New("invalid user name")
	ErrInvalidUserEmail  = errors.New("invalid user email")
	ErrUserNameTooLong   = errors.New("user name is too long")
	ErrUserEmailTooLong  = errors.New("user email is too long")
	
	// Role errors
	ErrRoleNotFound      = errors.New("role not found")
	ErrInvalidRoleID     = errors.New("invalid role ID")
	ErrInvalidRoleName   = errors.New("invalid role name")
	
	// Task type errors
	ErrTaskTypeNotFound              = errors.New("task type not found")
	ErrInvalidTaskTypeID             = errors.New("invalid task type ID")
	ErrInvalidTaskTypeName           = errors.New("invalid task type name")
	ErrTaskTypeNameTooLong           = errors.New("task type name is too long")
	ErrTaskTypeDescriptionTooLong    = errors.New("task type description is too long")
	
	// Machine errors
	ErrMachineNotFound               = errors.New("machine not found")
	ErrInvalidMachineID              = errors.New("invalid machine ID")
	ErrInvalidMachineName            = errors.New("invalid machine name")
	ErrMachineNameTooLong            = errors.New("machine name is too long")
	ErrInvalidMachineQuantity        = errors.New("machine quantity must be greater than 0")
	ErrInvalidMachineEstimateTime    = errors.New("machine estimate time must be greater than 0")
	
	// User task assignment errors
	ErrInvalidUserRole               = errors.New("only members can be assigned task types")
	ErrInvalidTaskTypeAssignment     = errors.New("invalid task type assignment")
) 