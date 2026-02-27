package services

import (
	"context"
	"testing"

	"github.com/suprt/planica_bi/backend/internal/models"
)

// MockDirectRepositoryForMarketing implements DirectRepositoryInterface for marketing service testing
type MockDirectRepositoryForMarketing struct {
	GetTotalsMonthlyFunc func(ctx context.Context, projectID uint, year int, month int) (*models.DirectTotalsMonthly, error)
}

func (m *MockDirectRepositoryForMarketing) CreateAccount(ctx context.Context, account *models.DirectAccount) error {
	return nil
}

func (m *MockDirectRepositoryForMarketing) GetAccountsByProjectID(ctx context.Context, projectID uint) ([]*models.DirectAccount, error) {
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) GetAccountByClientLogin(ctx context.Context, projectID uint, clientLogin string) (*models.DirectAccount, error) {
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) GetAccountByID(ctx context.Context, id uint) (*models.DirectAccount, error) {
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) UpdateAccount(ctx context.Context, account *models.DirectAccount) error {
	return nil
}

func (m *MockDirectRepositoryForMarketing) DeleteAccount(ctx context.Context, id uint) error {
	return nil
}

func (m *MockDirectRepositoryForMarketing) GetCampaignsByProjectID(ctx context.Context, projectID uint) ([]*models.DirectCampaign, error) {
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) GetCampaignsByAccountID(ctx context.Context, accountID uint) ([]*models.DirectCampaign, error) {
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) CreateCampaign(ctx context.Context, campaign *models.DirectCampaign) error {
	return nil
}

func (m *MockDirectRepositoryForMarketing) GetCampaignByID(ctx context.Context, id uint) (*models.DirectCampaign, error) {
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) GetCampaignMonthly(ctx context.Context, projectID uint, year int, month int) ([]*models.DirectCampaignMonthly, error) {
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) GetCampaignMonthlyByCampaignID(ctx context.Context, projectID uint, campaignID uint, year int, month int) (*models.DirectCampaignMonthly, error) {
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) SaveCampaignMonthly(campaign *models.DirectCampaignMonthly) error {
	return nil
}

func (m *MockDirectRepositoryForMarketing) GetTotalsMonthly(ctx context.Context, projectID uint, year int, month int) (*models.DirectTotalsMonthly, error) {
	if m.GetTotalsMonthlyFunc != nil {
		return m.GetTotalsMonthlyFunc(ctx, projectID, year, month)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForMarketing) SaveTotalsMonthly(totals *models.DirectTotalsMonthly) error {
	return nil
}

func TestMarketingService_GetMarketingData(t *testing.T) {
	tests := []struct {
		name        string
		projectID   uint
		mockSetup   func() *MockDirectRepositoryForMarketing
		wantErr     bool
		wantErrText string
	}{
		{
			name:      "успешное получение маркетинговых данных",
			projectID: 1,
			mockSetup: func() *MockDirectRepositoryForMarketing {
				return &MockDirectRepositoryForMarketing{
					GetTotalsMonthlyFunc: func(ctx context.Context, projectID uint, year int, month int) (*models.DirectTotalsMonthly, error) {
						return &models.DirectTotalsMonthly{
							Clicks:      1000,
							Impressions: 50000,
							Cost:        10000.0,
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:      "нет данных за месяцы (возвращает пустую структуру)",
			projectID: 1,
			mockSetup: func() *MockDirectRepositoryForMarketing {
				return &MockDirectRepositoryForMarketing{
					GetTotalsMonthlyFunc: func(ctx context.Context, projectID uint, year int, month int) (*models.DirectTotalsMonthly, error) {
						return nil, nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewMarketingService(mockRepo)

			data, err := service.GetMarketingData(context.Background(), tt.projectID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if data == nil {
					t.Errorf("ожидался ответ с данными, но получили nil")
				}
			}
		})
	}
}
