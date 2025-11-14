package repositories

import (
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project
func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetAll retrieves all projects
func (r *ProjectRepository) GetAll() ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.Find(&projects).Error
	return projects, err
}

// Update updates a project
func (r *ProjectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

// Delete deletes a project
func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}
