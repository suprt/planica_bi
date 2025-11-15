package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
)

// ProjectService handles business logic for projects
type ProjectService struct {
	projectRepo *repositories.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

// generatePublicToken generates a unique public token for project
func (s *ProjectService) generatePublicToken() (string, error) {
	bytes := make([]byte, 32) // 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, project *models.Project) error {
	// Validate required fields
	if project.Name == "" {
		return errors.New("name is required")
	}
	if project.Slug == "" {
		return errors.New("slug is required")
	}

	// Generate public token if not provided
	if project.PublicToken == "" {
		token, err := s.generatePublicToken()
		if err != nil {
			return fmt.Errorf("failed to generate public token: %w", err)
		}
		project.PublicToken = token
	}

	return s.projectRepo.Create(ctx, project)
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(ctx context.Context, id uint) (*models.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

// GetAllProjects retrieves all projects for a user
// Admin gets all projects, others get only projects they have access to
func (s *ProjectService) GetAllProjects(ctx context.Context, userID uint, isAdmin bool) ([]*models.Project, error) {
	return s.projectRepo.GetByUserID(ctx, userID, isAdmin)
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(ctx context.Context, project *models.Project) error {
	// Validate required fields
	if project.Name == "" {
		return errors.New("name is required")
	}
	if project.Slug == "" {
		return errors.New("slug is required")
	}

	// Check if project exists
	_, err := s.projectRepo.GetByID(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	return s.projectRepo.Update(ctx, project)
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id uint) error {
	return s.projectRepo.Delete(ctx, id)
}

// GetProjectByPublicToken retrieves a project by public token
func (s *ProjectService) GetProjectByPublicToken(ctx context.Context, token string) (*models.Project, error) {
	return s.projectRepo.GetByPublicToken(ctx, token)
}
