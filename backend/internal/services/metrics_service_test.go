package services

import (
	"context"
	"testing"

	"github.com/suprt/planica_bi/backend/internal/models"
)

// MockMetricsRepository implements MetricsRepositoryInterface for metrics service testing
type MockMetricsRepository struct {
	CreateFunc                         func(ctx context.Context, metrics *models.MetricsMonthly) error
	SaveMonthlyMetricsFunc             func(metrics *models.MetricsMonthly) error
	GetByProjectIDFunc                 func(ctx context.Context, projectID uint, month string) (*models.MetricsMonthly, error)
	GetMonthlyMetricsFunc              func(ctx context.Context, projectID uint, year int, month int) (*models.MetricsMonthly, error)
	GetAgeMetricsFunc                  func(ctx context.Context, projectID uint, year int, month int) ([]*models.MetricsAgeMonthly, error)
	GetAgeMetricsByGroupFunc           func(ctx context.Context, projectID uint, year int, month int, ageGroup string) (*models.MetricsAgeMonthly, error)
	SaveAgeMetricsFunc                 func(metrics *models.MetricsAgeMonthly) error
	GetByProjectIDAndAgeFunc           func(ctx context.Context, projectID uint, month string) ([]*models.MetricsAgeMonthly, error)
	UpdateFunc                         func(ctx context.Context, metrics *models.MetricsMonthly) error
	DeleteFunc                         func(ctx context.Context, id uint) error
	GetAllMonthlyMetricsForProjectFunc func(ctx context.Context, projectID uint) ([]*models.MetricsMonthly, error)
}

func (m *MockMetricsRepository) Create(ctx context.Context, metrics *models.MetricsMonthly) error {
	return nil
}

func (m *MockMetricsRepository) SaveMonthlyMetrics(metrics *models.MetricsMonthly) error {
	if m.SaveMonthlyMetricsFunc != nil {
		return m.SaveMonthlyMetricsFunc(metrics)
	}
	return nil
}

func (m *MockMetricsRepository) GetByProjectID(ctx context.Context, projectID uint, month string) (*models.MetricsMonthly, error) {
	if m.GetByProjectIDFunc != nil {
		return m.GetByProjectIDFunc(ctx, projectID, month)
	}
	return nil, nil
}

func (m *MockMetricsRepository) GetMonthlyMetrics(ctx context.Context, projectID uint, year int, month int) (*models.MetricsMonthly, error) {
	if m.GetMonthlyMetricsFunc != nil {
		return m.GetMonthlyMetricsFunc(ctx, projectID, year, month)
	}
	return nil, nil
}

func (m *MockMetricsRepository) GetAgeMetrics(ctx context.Context, projectID uint, year int, month int) ([]*models.MetricsAgeMonthly, error) {
	if m.GetAgeMetricsFunc != nil {
		return m.GetAgeMetricsFunc(ctx, projectID, year, month)
	}
	return nil, nil
}

func (m *MockMetricsRepository) GetAgeMetricsByGroup(ctx context.Context, projectID uint, year int, month int, ageGroup string) (*models.MetricsAgeMonthly, error) {
	if m.GetAgeMetricsByGroupFunc != nil {
		return m.GetAgeMetricsByGroupFunc(ctx, projectID, year, month, ageGroup)
	}
	return nil, nil
}

func (m *MockMetricsRepository) SaveAgeMetrics(metrics *models.MetricsAgeMonthly) error {
	if m.SaveAgeMetricsFunc != nil {
		return m.SaveAgeMetricsFunc(metrics)
	}
	return nil
}

func (m *MockMetricsRepository) GetByProjectIDAndAge(ctx context.Context, projectID uint, month string) ([]*models.MetricsAgeMonthly, error) {
	if m.GetByProjectIDAndAgeFunc != nil {
		return m.GetByProjectIDAndAgeFunc(ctx, projectID, month)
	}
	return nil, nil
}

func (m *MockMetricsRepository) Update(ctx context.Context, metrics *models.MetricsMonthly) error {
	return nil
}

func (m *MockMetricsRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m *MockMetricsRepository) GetAllMonthlyMetricsForProject(ctx context.Context, projectID uint) ([]*models.MetricsMonthly, error) {
	if m.GetAllMonthlyMetricsForProjectFunc != nil {
		return m.GetAllMonthlyMetricsForProjectFunc(ctx, projectID)
	}
	return nil, nil
}

func TestMetricsService_GetMetricsWithData(t *testing.T) {
	tests := []struct {
		name      string
		projectID uint
		mockSetup func() *MockMetricsRepository
		wantErr   bool
	}{
		{
			name:      "успешное получение метрик",
			projectID: 1,
			mockSetup: func() *MockMetricsRepository {
				return &MockMetricsRepository{
					GetMonthlyMetricsFunc: func(ctx context.Context, projectID uint, year int, month int) (*models.MetricsMonthly, error) {
						return &models.MetricsMonthly{
							Visits:     10000,
							Users:      5000,
							BounceRate: 35.5,
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:      "нет метрик (возвращает пустую структуру)",
			projectID: 1,
			mockSetup: func() *MockMetricsRepository {
				return &MockMetricsRepository{
					GetMonthlyMetricsFunc: func(ctx context.Context, projectID uint, year int, month int) (*models.MetricsMonthly, error) {
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
			service := NewMetricsService(mockRepo)

			metrics, err := service.GetMetricsWithData(context.Background(), tt.projectID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if metrics == nil {
					t.Errorf("ожидался ответ с данными, но получили nil")
				}
			}
		})
	}
}
