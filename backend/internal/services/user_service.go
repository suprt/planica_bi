package services

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
)

// UserService handles business logic for user management
type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUserRequest represents request to create a user
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}

// UpdateUserRequest represents request to update a user
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// UserResponse represents user data for API responses (without password)
type UserResponse struct {
	ID       uint                   `json:"id"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	IsActive bool                   `json:"is_active"`
	Projects []UserProjectResponse  `json:"projects,omitempty"`
}

// UserProjectResponse represents user's project and role
type UserProjectResponse struct {
	ProjectID   uint   `json:"project_id"`
	ProjectName string `json:"project_name"`
	Role        string `json:"role"`
}

// AssignRoleRequest represents request to assign role to user
type AssignRoleRequest struct {
	UserID    uint   `json:"user_id"`
	ProjectID uint   `json:"project_id"`
	Role      string `json:"role"`
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(ctx context.Context) ([]UserResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
		}
	}

	return response, nil
}

// GetUserByID retrieves user by ID with their projects
func (s *UserService) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user projects
	projectRoles, err := s.userRepo.GetUserProjects(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user projects: %w", err)
	}

	projects := make([]UserProjectResponse, len(projectRoles))
	for i, pr := range projectRoles {
		projects[i] = UserProjectResponse{
			ProjectID:   pr.ProjectID,
			ProjectName: pr.Project.Name,
			Role:        pr.Role,
		}
	}

	return &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
		Projects: projects,
	}, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		IsActive: req.IsActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, userID uint, req *UpdateUserRequest) (*UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		if *req.Name == "" {
			return nil, errors.New("name cannot be empty")
		}
		user.Name = *req.Name
	}

	if req.Email != nil {
		if *req.Email == "" {
			return nil, errors.New("email cannot be empty")
		}
		// Check if email is already taken by another user
		existing, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil && existing != nil && existing.ID != userID {
			return nil, errors.New("email already taken by another user")
		}
		user.Email = *req.Email
	}

	if req.Password != nil {
		if len(*req.Password) < 8 {
			return nil, errors.New("password must be at least 8 characters")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = string(hashedPassword)
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Save updates
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		IsActive: user.IsActive,
	}, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete user (cascade will handle project roles)
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetProjectUsers retrieves all users with roles for a project
func (s *UserService) GetProjectUsers(ctx context.Context, projectID uint) ([]UserResponse, error) {
	projectRoles, err := s.userRepo.GetProjectUsers(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project users: %w", err)
	}

	response := make([]UserResponse, len(projectRoles))
	for i, pr := range projectRoles {
		response[i] = UserResponse{
			ID:       pr.User.ID,
			Name:     pr.User.Name,
			Email:    pr.User.Email,
			IsActive: pr.User.IsActive,
			Projects: []UserProjectResponse{
				{
					ProjectID:   pr.ProjectID,
					Role:        pr.Role,
				},
			},
		}
	}

	return response, nil
}

// AssignRole assigns a role to user in project
func (s *UserService) AssignRole(ctx context.Context, req *AssignRoleRequest) error {
	// Validate role
	if req.Role != "admin" && req.Role != "manager" && req.Role != "client" {
		return errors.New("invalid role, must be 'admin', 'manager', or 'client'")
	}

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Assign role
	if err := s.userRepo.AssignRole(ctx, req.UserID, req.ProjectID, req.Role); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// UpdateRole updates user's role in project
func (s *UserService) UpdateRole(ctx context.Context, userID, projectID uint, role string) error {
	// Validate role
	if role != "admin" && role != "manager" && role != "client" {
		return errors.New("invalid role, must be 'admin', 'manager', or 'client'")
	}

	// Update role
	if err := s.userRepo.UpdateRole(ctx, userID, projectID, role); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

// RemoveRole removes user's role from project
func (s *UserService) RemoveRole(ctx context.Context, userID, projectID uint) error {
	// Remove role
	if err := s.userRepo.RemoveRole(ctx, userID, projectID); err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	return nil
}

