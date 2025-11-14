package services

import (
	"context"
	"errors"

	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/models"
	"gitlab.ugatu.su/gantseff/planica_bi/backend/internal/repositories"
	"gorm.io/gorm"
)

// DirectService handles business logic for Yandex.Direct accounts
type DirectService struct {
	directRepo *repositories.DirectRepository
}

// NewDirectService creates a new Direct service
func NewDirectService(directRepo *repositories.DirectRepository) *DirectService {
	return &DirectService{
		directRepo: directRepo,
	}
}

// CreateAccount creates a new Direct account
func (s *DirectService) CreateAccount(ctx context.Context, account *models.DirectAccount) error {
	// Validate required fields
	if account.ProjectID == 0 {
		return errors.New("project_id is required")
	}
	if account.ClientLogin == "" {
		return errors.New("client_login is required")
	}

	// Check if account with this ClientLogin already exists for this project
	existing, err := s.directRepo.GetAccountByClientLogin(ctx, account.ProjectID, account.ClientLogin)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		return errors.New("account with this ClientLogin already exists for this project")
	}

	return s.directRepo.CreateAccount(ctx, account)
}

// GetAccountsByProject retrieves all Direct accounts for a project
func (s *DirectService) GetAccountsByProject(ctx context.Context, projectID uint) ([]*models.DirectAccount, error) {
	return s.directRepo.GetAccountsByProjectID(ctx, projectID)
}
