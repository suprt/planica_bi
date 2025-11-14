package services

import "gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"

// ProjectService handles business logic for projects
type ProjectService struct {
	// TODO: add repositories
}

// NewProjectService creates a new project service
func NewProjectService() *ProjectService {
	return &ProjectService{}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(project *models.Project) error {
	// TODO: implement
	return nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(id uint) (*models.Project, error) {
	// TODO: implement
	return nil, nil
}

// GetAllProjects retrieves all projects
func (s *ProjectService) GetAllProjects() ([]*models.Project, error) {
	// TODO: implement
	return nil, nil
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(project *models.Project) error {
	// TODO: implement
	return nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id uint) error {
	// TODO: implement
	return nil
}

