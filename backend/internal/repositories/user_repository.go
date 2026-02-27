package repositories

import (
	"context"

	"github.com/suprt/planica_bi/backend/internal/middleware"
	"github.com/suprt/planica_bi/backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("ProjectRoles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin updates user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", gorm.Expr("NOW()")).Error
}

// GetUserProjectRole retrieves user's role for a specific project
func (r *UserRepository) GetUserProjectRole(ctx context.Context, userID, projectID uint) (*models.UserProjectRole, error) {
	var role models.UserProjectRole
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetUserProjects retrieves all projects for a user with their roles
func (r *UserRepository) GetUserProjects(ctx context.Context, userID uint) ([]models.UserProjectRole, error) {
	var roles []models.UserProjectRole
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Project").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// HasProjectAccess checks if user has access to a project
func (r *UserRepository) HasProjectAccess(ctx context.Context, userID, projectID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.UserProjectRole{}).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsAdmin checks if user is admin (has admin role on any project)
func (r *UserRepository) IsAdmin(ctx context.Context, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.UserProjectRole{}).
		Where("user_id = ? AND role = ?", userID, "admin").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Preload("ProjectRoles").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetAllPaginated retrieves paginated users
func (r *UserRepository) GetAllPaginated(ctx context.Context, pagination *middleware.Pagination) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Preload("ProjectRoles").
		Order(pagination.Sort + " " + pagination.Order).
		Limit(pagination.PerPage).
		Offset(pagination.Offset).
		Find(&users).Error

	return users, total, err
}

// Update updates user information
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, userID).Error
}

// GetProjectUsers retrieves all users for a project with their roles
func (r *UserRepository) GetProjectUsers(ctx context.Context, projectID uint) ([]models.UserProjectRole, error) {
	var roles []models.UserProjectRole
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Preload("User").
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// AssignRole assigns a role to user in project
// If role already exists, it will be updated
func (r *UserRepository) AssignRole(ctx context.Context, userID, projectID uint, role string) error {
	// Check if role already exists
	var existing models.UserProjectRole
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		First(&existing).Error

	if err == nil {
		// Role exists, update it
		existing.Role = role
		return r.db.WithContext(ctx).Save(&existing).Error
	}

	// Role doesn't exist, create new
	newRole := &models.UserProjectRole{
		UserID:    userID,
		ProjectID: projectID,
		Role:      role,
	}
	return r.db.WithContext(ctx).Create(newRole).Error
}

// UpdateRole updates user's role in project
func (r *UserRepository) UpdateRole(ctx context.Context, userID, projectID uint, role string) error {
	return r.db.WithContext(ctx).Model(&models.UserProjectRole{}).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Update("role", role).Error
}

// RemoveRole removes user's role from project
func (r *UserRepository) RemoveRole(ctx context.Context, userID, projectID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Delete(&models.UserProjectRole{}).Error
}
