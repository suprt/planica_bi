package services

import (
	"context"

	"github.com/suprt/planica_bi/backend/internal/middleware"
	"github.com/suprt/planica_bi/backend/internal/models"
)

// ============================================================================
// Repository Interfaces (defined by services - the consumer)
// ============================================================================

// UserRepositoryInterface defines methods for user data access
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetAllPaginated(ctx context.Context, pagination *middleware.Pagination) ([]models.User, int64, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userID uint) error
	UpdateLastLogin(ctx context.Context, userID uint) error
	GetUserProjectRole(ctx context.Context, userID, projectID uint) (*models.UserProjectRole, error)
	GetUserProjects(ctx context.Context, userID uint) ([]models.UserProjectRole, error)
	GetProjectUsers(ctx context.Context, projectID uint) ([]models.UserProjectRole, error)
	AssignRole(ctx context.Context, userID, projectID uint, role string) error
	UpdateRole(ctx context.Context, userID, projectID uint, role string) error
	RemoveRole(ctx context.Context, userID, projectID uint) error
	HasProjectAccess(ctx context.Context, userID, projectID uint) (bool, error)
	IsAdmin(ctx context.Context, userID uint) (bool, error)
}

// ProjectRepositoryInterface defines methods for project data access
type ProjectRepositoryInterface interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id uint) (*models.Project, error)
	GetAll(ctx context.Context) ([]*models.Project, error)
	GetAllPaginated(ctx context.Context, pagination *middleware.Pagination) ([]*models.Project, int64, error)
	GetByUserID(ctx context.Context, userID uint, isAdmin bool) ([]*models.Project, error)
	GetByUserIDPaginated(ctx context.Context, userID uint, isAdmin bool, pagination *middleware.Pagination) ([]*models.Project, int64, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id uint) error
	GetByPublicToken(ctx context.Context, token string) (*models.Project, error)
}

// CounterRepositoryInterface defines methods for counter data access
type CounterRepositoryInterface interface {
	Create(ctx context.Context, counter *models.YandexCounter) error
	GetByProjectID(ctx context.Context, projectID uint) ([]*models.YandexCounter, error)
	GetByID(ctx context.Context, id uint) (*models.YandexCounter, error)
	GetByCounterID(ctx context.Context, counterID int64) (*models.YandexCounter, error)
	Update(ctx context.Context, counter *models.YandexCounter) error
}

// DirectRepositoryInterface defines methods for Direct data access
type DirectRepositoryInterface interface {
	CreateAccount(ctx context.Context, account *models.DirectAccount) error
	GetAccountsByProjectID(ctx context.Context, projectID uint) ([]*models.DirectAccount, error)
	GetAccountByClientLogin(ctx context.Context, projectID uint, clientLogin string) (*models.DirectAccount, error)
	GetCampaignsByAccountID(ctx context.Context, accountID uint) ([]*models.DirectCampaign, error)
	CreateCampaign(ctx context.Context, campaign *models.DirectCampaign) error
	GetCampaignByID(ctx context.Context, id uint) (*models.DirectCampaign, error)
	GetCampaignsByProjectID(ctx context.Context, projectID uint) ([]*models.DirectCampaign, error)
	GetCampaignMonthly(ctx context.Context, projectID uint, year int, month int) ([]*models.DirectCampaignMonthly, error)
	GetCampaignMonthlyByCampaignID(ctx context.Context, projectID uint, directCampaignID uint, year int, month int) (*models.DirectCampaignMonthly, error)
	SaveCampaignMonthly(metrics *models.DirectCampaignMonthly) error
	GetTotalsMonthly(ctx context.Context, projectID uint, year int, month int) (*models.DirectTotalsMonthly, error)
	SaveTotalsMonthly(totals *models.DirectTotalsMonthly) error
}

// GoalRepositoryInterface defines methods for goal data access
type GoalRepositoryInterface interface {
	Create(ctx context.Context, goal *models.Goal) error
	GetByCounterID(ctx context.Context, counterID uint) ([]*models.Goal, error)
	GetByCounterIDs(ctx context.Context, counterIDs []uint) ([]*models.Goal, error)
	GetByID(ctx context.Context, id uint) (*models.Goal, error)
	GetByGoalID(ctx context.Context, counterID uint, goalID int64) (*models.Goal, error)
	Update(ctx context.Context, goal *models.Goal) error
	Delete(ctx context.Context, id uint) error
	DeleteByCounterID(ctx context.Context, counterID uint) error
}

// MetricsRepositoryInterface defines methods for metrics data access
type MetricsRepositoryInterface interface {
	GetMonthlyMetrics(ctx context.Context, projectID uint, year int, month int) (*models.MetricsMonthly, error)
	SaveMonthlyMetrics(metrics *models.MetricsMonthly) error
	GetAgeMetrics(ctx context.Context, projectID uint, year int, month int) ([]*models.MetricsAgeMonthly, error)
	SaveAgeMetrics(metrics *models.MetricsAgeMonthly) error
	GetAgeMetricsByGroup(ctx context.Context, projectID uint, year int, month int, ageGroup string) (*models.MetricsAgeMonthly, error)
	GetAllMonthlyMetricsForProject(ctx context.Context, projectID uint) ([]*models.MetricsMonthly, error)
}

// SEORepositoryInterface defines methods for SEO data access
type SEORepositoryInterface interface {
	GetSEOQueries(ctx context.Context, projectID uint, year int, month int) ([]*models.SEOQueriesMonthly, error)
}
