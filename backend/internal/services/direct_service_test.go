package services

import (
	"context"
	"testing"

	"github.com/suprt/planica_bi/backend/internal/models"
)

// MockDirectRepositoryForDirectService implements DirectRepositoryInterface for direct service testing
type MockDirectRepositoryForDirectService struct {
	CreateAccountFunc                  func(ctx context.Context, account *models.DirectAccount) error
	GetAccountsByProjectIDFunc         func(ctx context.Context, projectID uint) ([]*models.DirectAccount, error)
	GetAccountByClientLoginFunc        func(ctx context.Context, projectID uint, clientLogin string) (*models.DirectAccount, error)
	GetAccountByIDFunc                 func(ctx context.Context, id uint) (*models.DirectAccount, error)
	UpdateAccountFunc                  func(ctx context.Context, account *models.DirectAccount) error
	DeleteAccountFunc                  func(ctx context.Context, id uint) error
	GetCampaignsByProjectIDFunc        func(ctx context.Context, projectID uint) ([]*models.DirectCampaign, error)
	GetCampaignsByAccountIDFunc        func(ctx context.Context, accountID uint) ([]*models.DirectCampaign, error)
	CreateCampaignFunc                 func(ctx context.Context, campaign *models.DirectCampaign) error
	GetCampaignByIDFunc                func(ctx context.Context, id uint) (*models.DirectCampaign, error)
	GetCampaignMonthlyFunc             func(ctx context.Context, projectID uint, year int, month int) ([]*models.DirectCampaignMonthly, error)
	GetCampaignMonthlyByCampaignIDFunc func(ctx context.Context, projectID uint, campaignID uint, year int, month int) (*models.DirectCampaignMonthly, error)
	SaveCampaignMonthlyFunc            func(campaign *models.DirectCampaignMonthly) error
	GetTotalsMonthlyFunc               func(ctx context.Context, projectID uint, year int, month int) (*models.DirectTotalsMonthly, error)
	SaveTotalsMonthlyFunc              func(totals *models.DirectTotalsMonthly) error
}

func (m *MockDirectRepositoryForDirectService) CreateAccount(ctx context.Context, account *models.DirectAccount) error {
	if m.CreateAccountFunc != nil {
		return m.CreateAccountFunc(ctx, account)
	}
	return nil
}

func (m *MockDirectRepositoryForDirectService) GetAccountsByProjectID(ctx context.Context, projectID uint) ([]*models.DirectAccount, error) {
	if m.GetAccountsByProjectIDFunc != nil {
		return m.GetAccountsByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) GetAccountByClientLogin(ctx context.Context, projectID uint, clientLogin string) (*models.DirectAccount, error) {
	if m.GetAccountByClientLoginFunc != nil {
		return m.GetAccountByClientLoginFunc(ctx, projectID, clientLogin)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) GetAccountByID(ctx context.Context, id uint) (*models.DirectAccount, error) {
	if m.GetAccountByIDFunc != nil {
		return m.GetAccountByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) UpdateAccount(ctx context.Context, account *models.DirectAccount) error {
	if m.UpdateAccountFunc != nil {
		return m.UpdateAccountFunc(ctx, account)
	}
	return nil
}

func (m *MockDirectRepositoryForDirectService) DeleteAccount(ctx context.Context, id uint) error {
	if m.DeleteAccountFunc != nil {
		return m.DeleteAccountFunc(ctx, id)
	}
	return nil
}

func (m *MockDirectRepositoryForDirectService) GetCampaignsByProjectID(ctx context.Context, projectID uint) ([]*models.DirectCampaign, error) {
	if m.GetCampaignsByProjectIDFunc != nil {
		return m.GetCampaignsByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) GetCampaignsByAccountID(ctx context.Context, accountID uint) ([]*models.DirectCampaign, error) {
	if m.GetCampaignsByAccountIDFunc != nil {
		return m.GetCampaignsByAccountIDFunc(ctx, accountID)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) CreateCampaign(ctx context.Context, campaign *models.DirectCampaign) error {
	if m.CreateCampaignFunc != nil {
		return m.CreateCampaignFunc(ctx, campaign)
	}
	return nil
}

func (m *MockDirectRepositoryForDirectService) GetCampaignByID(ctx context.Context, id uint) (*models.DirectCampaign, error) {
	if m.GetCampaignByIDFunc != nil {
		return m.GetCampaignByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) GetCampaignMonthly(ctx context.Context, projectID uint, year int, month int) ([]*models.DirectCampaignMonthly, error) {
	if m.GetCampaignMonthlyFunc != nil {
		return m.GetCampaignMonthlyFunc(ctx, projectID, year, month)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) GetCampaignMonthlyByCampaignID(ctx context.Context, projectID uint, campaignID uint, year int, month int) (*models.DirectCampaignMonthly, error) {
	if m.GetCampaignMonthlyByCampaignIDFunc != nil {
		return m.GetCampaignMonthlyByCampaignIDFunc(ctx, projectID, campaignID, year, month)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) SaveCampaignMonthly(campaign *models.DirectCampaignMonthly) error {
	if m.SaveCampaignMonthlyFunc != nil {
		return m.SaveCampaignMonthlyFunc(campaign)
	}
	return nil
}

func (m *MockDirectRepositoryForDirectService) GetTotalsMonthly(ctx context.Context, projectID uint, year int, month int) (*models.DirectTotalsMonthly, error) {
	if m.GetTotalsMonthlyFunc != nil {
		return m.GetTotalsMonthlyFunc(ctx, projectID, year, month)
	}
	return nil, nil
}

func (m *MockDirectRepositoryForDirectService) SaveTotalsMonthly(totals *models.DirectTotalsMonthly) error {
	if m.SaveTotalsMonthlyFunc != nil {
		return m.SaveTotalsMonthlyFunc(totals)
	}
	return nil
}

func TestDirectService_CreateAccount(t *testing.T) {
	tests := []struct {
		name        string
		account     *models.DirectAccount
		mockSetup   func() *MockDirectRepositoryForDirectService
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешное создание аккаунта",
			account: &models.DirectAccount{
				ProjectID:   1,
				ClientLogin: "test_login",
			},
			mockSetup: func() *MockDirectRepositoryForDirectService {
				return &MockDirectRepositoryForDirectService{
					GetAccountByClientLoginFunc: func(ctx context.Context, projectID uint, clientLogin string) (*models.DirectAccount, error) {
						return nil, nil
					},
					CreateAccountFunc: func(ctx context.Context, account *models.DirectAccount) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "пустой project_id",
			account: &models.DirectAccount{
				ProjectID:   0,
				ClientLogin: "test_login",
			},
			mockSetup: func() *MockDirectRepositoryForDirectService {
				return &MockDirectRepositoryForDirectService{}
			},
			wantErr:     true,
			wantErrText: "project_id is required",
		},
		{
			name: "пустой client_login",
			account: &models.DirectAccount{
				ProjectID:   1,
				ClientLogin: "",
			},
			mockSetup: func() *MockDirectRepositoryForDirectService {
				return &MockDirectRepositoryForDirectService{}
			},
			wantErr:     true,
			wantErrText: "client_login is required",
		},
		{
			name: "аккаунт уже существует",
			account: &models.DirectAccount{
				ProjectID:   1,
				ClientLogin: "test_login",
			},
			mockSetup: func() *MockDirectRepositoryForDirectService {
				return &MockDirectRepositoryForDirectService{
					GetAccountByClientLoginFunc: func(ctx context.Context, projectID uint, clientLogin string) (*models.DirectAccount, error) {
						return &models.DirectAccount{ClientLogin: clientLogin}, nil
					},
				}
			},
			wantErr:     true,
			wantErrText: "account with this ClientLogin already exists for this project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewDirectService(mockRepo)

			err := service.CreateAccount(context.Background(), tt.account)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
					return
				}
				if tt.wantErrText != "" && err.Error() != tt.wantErrText {
					t.Errorf("ожидалась ошибка '%s', но получили '%s'", tt.wantErrText, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
			}
		})
	}
}

func TestDirectService_GetAccountsByProject(t *testing.T) {
	tests := []struct {
		name      string
		projectID uint
		mockSetup func() *MockDirectRepositoryForDirectService
		wantCount int
		wantErr   bool
	}{
		{
			name:      "успешное получение аккаунтов",
			projectID: 1,
			mockSetup: func() *MockDirectRepositoryForDirectService {
				return &MockDirectRepositoryForDirectService{
					GetAccountsByProjectIDFunc: func(ctx context.Context, projectID uint) ([]*models.DirectAccount, error) {
						return []*models.DirectAccount{
							{ID: 1, ClientLogin: "login1"},
							{ID: 2, ClientLogin: "login2"},
						}, nil
					},
				}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "нет аккаунтов",
			projectID: 1,
			mockSetup: func() *MockDirectRepositoryForDirectService {
				return &MockDirectRepositoryForDirectService{
					GetAccountsByProjectIDFunc: func(ctx context.Context, projectID uint) ([]*models.DirectAccount, error) {
						return []*models.DirectAccount{}, nil
					},
				}
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewDirectService(mockRepo)

			accounts, err := service.GetAccountsByProject(context.Background(), tt.projectID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if len(accounts) != tt.wantCount {
					t.Errorf("ожидалось %d аккаунтов, но получили %d", tt.wantCount, len(accounts))
				}
			}
		})
	}
}
