package services

import (
	"context"
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

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, project *models.Project) error {
	// Validate required fields
	if project.Name == "" {
		return errors.New("name is required")
	}
	if project.Slug == "" {
		return errors.New("slug is required")
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

// GetAllActiveProjects retrieves all active projects (for system tasks like cron)
func (s *ProjectService) GetAllActiveProjects(ctx context.Context) ([]*models.Project, error) {
	allProjects, err := s.projectRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Filter only active projects
	var activeProjects []*models.Project
	for _, project := range allProjects {
		if project.IsActive {
			activeProjects = append(activeProjects, project)
		}
	}

	return activeProjects, nil
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
