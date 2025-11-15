package handlers

import (
	"context"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/services"
)

// UserServiceInterface defines methods for user service operations
type UserServiceInterface interface {
	GetAllUsers(ctx context.Context) ([]services.UserResponse, error)
	GetUserByID(ctx context.Context, userID uint) (*services.UserResponse, error)
	CreateUser(ctx context.Context, req *services.CreateUserRequest) (*services.UserResponse, error)
	UpdateUser(ctx context.Context, userID uint, req *services.UpdateUserRequest) (*services.UserResponse, error)
	DeleteUser(ctx context.Context, userID uint) error
	GetProjectUsers(ctx context.Context, projectID uint) ([]services.UserResponse, error)
	AssignRole(ctx context.Context, req *services.AssignRoleRequest) error
	UpdateRole(ctx context.Context, userID, projectID uint, role string) error
	RemoveRole(ctx context.Context, userID, projectID uint) error
}
