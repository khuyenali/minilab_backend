package service

import (
	"errors"
)

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
	ErrInvalidMachineIDs             = errors.New("one or more machine IDs are invalid")
	
	// Machine errors
	ErrMachineNotFound               = errors.New("machine not found")
	ErrInvalidMachineID              = errors.New("invalid machine ID")
	ErrInvalidMachineName            = errors.New("invalid machine name")
	ErrMachineNameTooLong            = errors.New("machine name is too long")
	ErrInvalidMachineQuantity        = errors.New("machine quantity must be greater than 0")
	ErrInvalidMachineEstimateTime    = errors.New("machine estimate time must be greater than 0")
	ErrMachineInvalidTaskTypeID      = errors.New("invalid task type ID")
	ErrMachineCannotRemoveTaskType   = errors.New("cannot remove task type from machine - can only change to another task type")
	
	// User task assignment errors
	ErrInvalidUserRole               = errors.New("only members can be assigned task types")
	ErrInvalidTaskTypeAssignment     = errors.New("invalid task type assignment")
	ErrInvalidTaskTypeIDs            = errors.New("one or more task type IDs are invalid")
	
	// Task errors
	ErrTaskNotFound                  = errors.New("task not found")
	ErrInvalidTaskID                 = errors.New("invalid task ID")
	ErrInvalidTaskName               = errors.New("invalid task name")
	ErrTaskNameTooLong               = errors.New("task name is too long")
	ErrTaskNoteTooLong               = errors.New("task note is too long")
	ErrInvalidTaskStatus             = errors.New("invalid task status")
	ErrInvalidTaskPriority           = errors.New("invalid task priority")
	ErrInvalidTaskStatusTransition   = errors.New("invalid task status transition")
	ErrTaskReportRequired            = errors.New("task report is required when finishing a task")
	
	// Sub-task errors
	ErrSubTaskNotFound               = errors.New("sub-task not found")
	ErrInvalidSubTaskID              = errors.New("invalid sub-task ID")
	ErrInvalidSubTaskName            = errors.New("invalid sub-task name")
	ErrSubTaskNameTooLong            = errors.New("sub-task name is too long")
	ErrSubTaskDescriptionTooLong     = errors.New("sub-task description is too long")
	ErrInvalidSubTaskStatus          = errors.New("invalid sub-task status")
	ErrInvalidSubTaskEstimateEffort  = errors.New("invalid sub-task estimate effort")
	
	// Assignment errors
	ErrAssignmentNotFound            = errors.New("assignment not found")
	ErrInvalidAssignmentID           = errors.New("invalid assignment ID")
	ErrAssignmentReportTooLong       = errors.New("assignment report is too long")
	ErrAssignmentReportRequired      = errors.New("assignment report is required")
	ErrAssignmentNotInPendingStatus  = errors.New("assignment must be in pending status to move to processing")
	ErrAssignmentNotInProcessingStatus = errors.New("assignment must be in processing status to move to finish")
	ErrInvalidSubTaskForAssignment   = errors.New("invalid sub-task ID for assignment")
	ErrInvalidUserForAssignment      = errors.New("invalid user ID for assignment")
	ErrAssignmentAlreadyExists       = errors.New("assignment already exists for this user and sub-task")
	ErrInvalidAssignmentStatus       = errors.New("invalid assignment status")
) 