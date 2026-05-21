package domain

import (
	"errors"
)

var (
	ErrValidation            = errors.New("validation failed")
	ErrEmailExists           = errors.New("user with this email already exists")
	ErrUsernameExists        = errors.New("user with this username already exists")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrUnauthorized          = errors.New("unauthorized access")
	ErrUserNotFound          = errors.New("user not found")
	ErrTaskNotFound          = errors.New("task not found")
	ErrEventNotFound         = errors.New("event not found")
	ErrCannotPromoteYourself = errors.New("cannot promote yourself")
	ErrCannotDemoteYourself  = errors.New("cannot demote yourself")
	ErrFailedToCreateTask    = errors.New("failed to create task")
	ErrFailedToCreateEvent   = errors.New("failed to create event")
	ErrFailedToGetUserTasks  = errors.New("failed to get user tasks")
	ErrFailedToGetUserEvents = errors.New("failed to get user events")
	ErrInvalidUUID           = errors.New("invalid UUID format")
	ErrForbidden             = errors.New("forbidden")

	ErrNoPermissionUpdateUser  = errors.New("you do not have the permissions to update this user")
	ErrNoPermissionDeleteUser  = errors.New("you do not have the permissions to delete this user")
	ErrNoPermissionViewTask    = errors.New("you don't have permission to view this task")
	ErrNoPermissionUpdateTask  = errors.New("you do not have the permissions to update this task")
	ErrNoPermissionDeleteTask  = errors.New("you do not have the permissions to delete this task")
	ErrNoPermissionViewEvent   = errors.New("you don't have permission to view this event")
	ErrNoPermissionUpdateEvent = errors.New("you don't have permission to update this event")
	ErrNoPermissionDeleteEvent = errors.New("you don't have permission to delete this event")
)
