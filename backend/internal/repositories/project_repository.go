package repositories

import "gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	// TODO: add database connection
}

// NewProjectRepository creates a new project repository
func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

// Create creates a new project
func (r *ProjectRepository) Create(project *models.Project) error {
	// TODO: implement
	return nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id uint) (*models.Project, error) {
	// TODO: implement
	return nil, nil
}

// GetAll retrieves all projects
func (r *ProjectRepository) GetAll() ([]*models.Project, error) {
	// TODO: implement
	return nil, nil
}

// Update updates a project
func (r *ProjectRepository) Update(project *models.Project) error {
	// TODO: implement
	return nil
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id uint) error {
	// TODO: implement
	return nil
}
