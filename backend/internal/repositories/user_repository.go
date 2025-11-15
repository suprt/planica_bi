package repositories

import (
	"context"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
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

