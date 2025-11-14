package services

import (
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
func (s *ProjectService) CreateProject(project *models.Project) error {
	return s.projectRepo.Create(project)
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(id uint) (*models.Project, error) {
	return s.projectRepo.GetByID(id)
}

// GetAllProjects retrieves all projects
func (s *ProjectService) GetAllProjects() ([]*models.Project, error) {
	return s.projectRepo.GetAll()
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(project *models.Project) error {
	return s.projectRepo.Update(project)
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id uint) error {
	return s.projectRepo.Delete(id)
}

