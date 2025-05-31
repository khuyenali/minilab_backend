package service

import "errors"

// Service errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserEmailExists   = errors.New("user with this email already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrInvalidUserName   = errors.New("invalid user name")
	ErrInvalidUserEmail  = errors.New("invalid user email")
	ErrUserNameTooLong   = errors.New("user name is too long")
	ErrUserEmailTooLong  = errors.New("user email is too long")
) 