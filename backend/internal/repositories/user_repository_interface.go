package repositories

import (
	"context"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
)

// UserRepositoryInterface defines methods for user repository operations
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
	GetUserProjectRole(ctx context.Context, userID, projectID uint) (*models.UserProjectRole, error)
	GetUserProjects(ctx context.Context, userID uint) ([]models.UserProjectRole, error)
	HasProjectAccess(ctx context.Context, userID, projectID uint) (bool, error)
	IsAdmin(ctx context.Context, userID uint) (bool, error)
}

