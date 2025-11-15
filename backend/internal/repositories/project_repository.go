package repositories

import (
	"context"

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
func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(ctx context.Context, id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetAll retrieves all projects
func (r *ProjectRepository) GetAll(ctx context.Context) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.WithContext(ctx).Find(&projects).Error
	return projects, err
}

// GetByUserID retrieves all projects for a specific user
// Admin gets all projects, others get only projects they have access to
func (r *ProjectRepository) GetByUserID(ctx context.Context, userID uint, isAdmin bool) ([]*models.Project, error) {
	var projects []*models.Project

	if isAdmin {
		// Admin sees all projects
		return r.GetAll(ctx)
	}

	// Get projects through user_project_roles
	err := r.db.WithContext(ctx).
		Table("projects").
		Joins("INNER JOIN user_project_roles ON projects.id = user_project_roles.project_id").
		Where("user_project_roles.user_id = ?", userID).
		Find(&projects).Error

	return projects, err
}

// Update updates a project
func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

// Delete deletes a project
func (r *ProjectRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, id).Error
}
