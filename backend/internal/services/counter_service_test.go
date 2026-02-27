package services

import (
	"context"
	"testing"

	"github.com/suprt/planica_bi/backend/internal/models"
)

// MockCounterRepository implements CounterRepositoryInterface for counter service testing
type MockCounterRepository struct {
	CreateFunc         func(ctx context.Context, counter *models.YandexCounter) error
	GetByProjectIDFunc func(ctx context.Context, projectID uint) ([]*models.YandexCounter, error)
	GetByIDFunc        func(ctx context.Context, id uint) (*models.YandexCounter, error)
	GetByCounterIDFunc func(ctx context.Context, counterID int64) (*models.YandexCounter, error)
	UpdateFunc         func(ctx context.Context, counter *models.YandexCounter) error
	DeleteFunc         func(ctx context.Context, id uint) error
}

func (m *MockCounterRepository) Create(ctx context.Context, counter *models.YandexCounter) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, counter)
	}
	return nil
}

func (m *MockCounterRepository) GetByProjectID(ctx context.Context, projectID uint) ([]*models.YandexCounter, error) {
	if m.GetByProjectIDFunc != nil {
		return m.GetByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockCounterRepository) GetByID(ctx context.Context, id uint) (*models.YandexCounter, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCounterRepository) GetByCounterID(ctx context.Context, counterID int64) (*models.YandexCounter, error) {
	if m.GetByCounterIDFunc != nil {
		return m.GetByCounterIDFunc(ctx, counterID)
	}
	return nil, nil
}

func (m *MockCounterRepository) Update(ctx context.Context, counter *models.YandexCounter) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, counter)
	}
	return nil
}

func (m *MockCounterRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestCounterService_CreateCounter(t *testing.T) {
	tests := []struct {
		name        string
		counter     *models.YandexCounter
		mockSetup   func() *MockCounterRepository
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешное создание счётчика",
			counter: &models.YandexCounter{
				ProjectID: 1,
				CounterID: 12345,
				Name:      "Test Counter",
				IsPrimary: false,
			},
			mockSetup: func() *MockCounterRepository {
				return &MockCounterRepository{
					GetByCounterIDFunc: func(ctx context.Context, counterID int64) (*models.YandexCounter, error) {
						return nil, nil
					},
					CreateFunc: func(ctx context.Context, counter *models.YandexCounter) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "пустой project_id",
			counter: &models.YandexCounter{
				ProjectID: 0,
				CounterID: 12345,
			},
			mockSetup: func() *MockCounterRepository {
				return &MockCounterRepository{}
			},
			wantErr:     true,
			wantErrText: "project_id is required",
		},
		{
			name: "пустой counter_id",
			counter: &models.YandexCounter{
				ProjectID: 1,
				CounterID: 0,
			},
			mockSetup: func() *MockCounterRepository {
				return &MockCounterRepository{}
			},
			wantErr:     true,
			wantErrText: "counter_id is required",
		},
		{
			name: "счётчик уже существует",
			counter: &models.YandexCounter{
				ProjectID: 1,
				CounterID: 12345,
			},
			mockSetup: func() *MockCounterRepository {
				return &MockCounterRepository{
					GetByCounterIDFunc: func(ctx context.Context, counterID int64) (*models.YandexCounter, error) {
						return &models.YandexCounter{CounterID: counterID}, nil
					},
				}
			},
			wantErr:     true,
			wantErrText: "counter with this CounterID already exists",
		},
		{
			name: "установка primary сбрасывает другие primary",
			counter: &models.YandexCounter{
				ProjectID: 1,
				CounterID: 12345,
				IsPrimary: true,
			},
			mockSetup: func() *MockCounterRepository {
				return &MockCounterRepository{
					GetByCounterIDFunc: func(ctx context.Context, counterID int64) (*models.YandexCounter, error) {
						return nil, nil
					},
					GetByProjectIDFunc: func(ctx context.Context, projectID uint) ([]*models.YandexCounter, error) {
						return []*models.YandexCounter{
							{ID: 1, ProjectID: projectID, IsPrimary: true},
						}, nil
					},
					UpdateFunc: func(ctx context.Context, counter *models.YandexCounter) error {
						return nil
					},
					CreateFunc: func(ctx context.Context, counter *models.YandexCounter) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewCounterService(mockRepo)

			err := service.CreateCounter(context.Background(), tt.counter)

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

func TestCounterService_GetCountersByProject(t *testing.T) {
	tests := []struct {
		name      string
		projectID uint
		mockSetup func() *MockCounterRepository
		wantCount int
		wantErr   bool
	}{
		{
			name:      "успешное получение счётчиков",
			projectID: 1,
			mockSetup: func() *MockCounterRepository {
				return &MockCounterRepository{
					GetByProjectIDFunc: func(ctx context.Context, projectID uint) ([]*models.YandexCounter, error) {
						return []*models.YandexCounter{
							{ID: 1, CounterID: 12345},
							{ID: 2, CounterID: 67890},
						}, nil
					},
				}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "нет счётчиков",
			projectID: 1,
			mockSetup: func() *MockCounterRepository {
				return &MockCounterRepository{
					GetByProjectIDFunc: func(ctx context.Context, projectID uint) ([]*models.YandexCounter, error) {
						return []*models.YandexCounter{}, nil
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
			service := NewCounterService(mockRepo)

			counters, err := service.GetCountersByProject(context.Background(), tt.projectID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if len(counters) != tt.wantCount {
					t.Errorf("ожидалось %d счётчиков, но получили %d", tt.wantCount, len(counters))
				}
			}
		})
	}
}
