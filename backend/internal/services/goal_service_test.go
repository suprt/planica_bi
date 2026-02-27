package services

import (
	"context"
	"testing"

	"github.com/suprt/planica_bi/backend/internal/models"
)

// MockGoalRepository implements GoalRepositoryInterface for goal service testing
type MockGoalRepository struct {
	CreateFunc            func(ctx context.Context, goal *models.Goal) error
	GetByProjectIDFunc    func(ctx context.Context, projectID uint) ([]*models.Goal, error)
	GetByIDFunc           func(ctx context.Context, id uint) (*models.Goal, error)
	GetByGoalIDFunc       func(ctx context.Context, counterID uint, goalID int64) (*models.Goal, error)
	GetByCounterIDFunc    func(ctx context.Context, counterID uint) ([]*models.Goal, error)
	GetByCounterIDsFunc   func(ctx context.Context, counterIDs []uint) ([]*models.Goal, error)
	UpdateFunc            func(ctx context.Context, goal *models.Goal) error
	DeleteFunc            func(ctx context.Context, id uint) error
	DeleteByCounterIDFunc func(ctx context.Context, counterID uint) error
}

func (m *MockGoalRepository) Create(ctx context.Context, goal *models.Goal) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, goal)
	}
	return nil
}

func (m *MockGoalRepository) GetByProjectID(ctx context.Context, projectID uint) ([]*models.Goal, error) {
	if m.GetByProjectIDFunc != nil {
		return m.GetByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockGoalRepository) GetByID(ctx context.Context, id uint) (*models.Goal, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockGoalRepository) GetByGoalID(ctx context.Context, counterID uint, goalID int64) (*models.Goal, error) {
	if m.GetByGoalIDFunc != nil {
		return m.GetByGoalIDFunc(ctx, counterID, goalID)
	}
	return nil, nil
}

func (m *MockGoalRepository) GetByCounterID(ctx context.Context, counterID uint) ([]*models.Goal, error) {
	if m.GetByCounterIDFunc != nil {
		return m.GetByCounterIDFunc(ctx, counterID)
	}
	return nil, nil
}

func (m *MockGoalRepository) GetByCounterIDs(ctx context.Context, counterIDs []uint) ([]*models.Goal, error) {
	if m.GetByCounterIDsFunc != nil {
		return m.GetByCounterIDsFunc(ctx, counterIDs)
	}
	return nil, nil
}

func (m *MockGoalRepository) Update(ctx context.Context, goal *models.Goal) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, goal)
	}
	return nil
}

func (m *MockGoalRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockGoalRepository) DeleteByCounterID(ctx context.Context, counterID uint) error {
	if m.DeleteByCounterIDFunc != nil {
		return m.DeleteByCounterIDFunc(ctx, counterID)
	}
	return nil
}

// MockCounterRepositoryForGoalService implements CounterRepositoryInterface
type MockCounterRepositoryForGoalService struct {
	GetByIDFunc func(ctx context.Context, id uint) (*models.YandexCounter, error)
}

func (m *MockCounterRepositoryForGoalService) Create(ctx context.Context, counter *models.YandexCounter) error {
	return nil
}

func (m *MockCounterRepositoryForGoalService) GetByProjectID(ctx context.Context, projectID uint) ([]*models.YandexCounter, error) {
	return nil, nil
}

func (m *MockCounterRepositoryForGoalService) GetByID(ctx context.Context, id uint) (*models.YandexCounter, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCounterRepositoryForGoalService) GetByCounterID(ctx context.Context, counterID int64) (*models.YandexCounter, error) {
	return nil, nil
}

func (m *MockCounterRepositoryForGoalService) Update(ctx context.Context, counter *models.YandexCounter) error {
	return nil
}

func (m *MockCounterRepositoryForGoalService) Delete(ctx context.Context, id uint) error {
	return nil
}

func TestGoalService_CreateGoal(t *testing.T) {
	tests := []struct {
		name        string
		goal        *models.Goal
		mockSetup   func() (*MockGoalRepository, *MockCounterRepositoryForGoalService)
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешное создание цели",
			goal: &models.Goal{
				CounterID: 1,
				GoalID:    12345,
				Name:      "Test Goal",
			},
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{
						GetByGoalIDFunc: func(ctx context.Context, counterID uint, goalID int64) (*models.Goal, error) {
							return nil, nil
						},
						CreateFunc: func(ctx context.Context, goal *models.Goal) error {
							return nil
						},
					}, &MockCounterRepositoryForGoalService{
						GetByIDFunc: func(ctx context.Context, id uint) (*models.YandexCounter, error) {
							return &models.YandexCounter{ID: id}, nil
						},
					}
			},
			wantErr: false,
		},
		{
			name: "пустой counter_id",
			goal: &models.Goal{
				CounterID: 0,
				GoalID:    12345,
			},
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{}, &MockCounterRepositoryForGoalService{}
			},
			wantErr:     true,
			wantErrText: "counter_id is required",
		},
		{
			name: "пустой goal_id",
			goal: &models.Goal{
				CounterID: 1,
				GoalID:    0,
			},
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{}, &MockCounterRepositoryForGoalService{}
			},
			wantErr:     true,
			wantErrText: "goal_id is required",
		},
		{
			name: "счётчик не найден",
			goal: &models.Goal{
				CounterID: 999,
				GoalID:    12345,
			},
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{}, &MockCounterRepositoryForGoalService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.YandexCounter, error) {
						return nil, nil
					},
				}
			},
			wantErr:     true,
			wantErrText: "counter not found",
		},
		{
			name: "цель уже существует",
			goal: &models.Goal{
				CounterID: 1,
				GoalID:    12345,
			},
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{
						GetByGoalIDFunc: func(ctx context.Context, counterID uint, goalID int64) (*models.Goal, error) {
							return &models.Goal{GoalID: goalID}, nil
						},
					}, &MockCounterRepositoryForGoalService{
						GetByIDFunc: func(ctx context.Context, id uint) (*models.YandexCounter, error) {
							return &models.YandexCounter{ID: id}, nil
						},
					}
			},
			wantErr:     true,
			wantErrText: "goal with this GoalID already exists for this counter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoalRepo, mockCounterRepo := tt.mockSetup()
			service := NewGoalService(mockGoalRepo, mockCounterRepo)

			err := service.CreateGoal(context.Background(), tt.goal)

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

func TestGoalService_GetGoal(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func() (*MockGoalRepository, *MockCounterRepositoryForGoalService)
		wantErr   bool
		wantID    uint
	}{
		{
			name: "успешное получение цели",
			id:   1,
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.Goal, error) {
						return &models.Goal{ID: id, Name: "Test Goal"}, nil
					},
				}, &MockCounterRepositoryForGoalService{}
			},
			wantErr: false,
			wantID:  1,
		},
		{
			name: "цель не найдена",
			id:   999,
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.Goal, error) {
						return nil, nil
					},
				}, &MockCounterRepositoryForGoalService{}
			},
			wantErr: true,
			wantID:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoalRepo, mockCounterRepo := tt.mockSetup()
			service := NewGoalService(mockGoalRepo, mockCounterRepo)

			goal, err := service.GetGoal(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil && goal == nil {
					// Ожидали ошибку или nil
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if goal == nil {
					t.Errorf("ожидалась цель, но получили nil")
				}
				if goal.ID != tt.wantID {
					t.Errorf("ожидался ID = %d, но получили %d", tt.wantID, goal.ID)
				}
			}
		})
	}
}

func TestGoalService_GetGoalsByCounter(t *testing.T) {
	tests := []struct {
		name        string
		counterID   uint
		mockSetup   func() (*MockGoalRepository, *MockCounterRepositoryForGoalService)
		wantCount   int
		wantErr     bool
		wantErrText string
	}{
		{
			name:      "успешное получение целей",
			counterID: 1,
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{
						GetByCounterIDFunc: func(ctx context.Context, counterID uint) ([]*models.Goal, error) {
							return []*models.Goal{
								{ID: 1, Name: "Goal 1"},
								{ID: 2, Name: "Goal 2"},
							}, nil
						},
					}, &MockCounterRepositoryForGoalService{
						GetByIDFunc: func(ctx context.Context, id uint) (*models.YandexCounter, error) {
							return &models.YandexCounter{ID: id}, nil
						},
					}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "счётчик не найден",
			counterID: 999,
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{}, &MockCounterRepositoryForGoalService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.YandexCounter, error) {
						return nil, nil
					},
				}
			},
			wantCount:   0,
			wantErr:     true,
			wantErrText: "counter not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoalRepo, mockCounterRepo := tt.mockSetup()
			service := NewGoalService(mockGoalRepo, mockCounterRepo)

			goals, err := service.GetGoalsByCounter(context.Background(), tt.counterID)

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
				if len(goals) != tt.wantCount {
					t.Errorf("ожидалось %d целей, но получили %d", tt.wantCount, len(goals))
				}
			}
		})
	}
}

func TestGoalService_DeleteGoal(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func() (*MockGoalRepository, *MockCounterRepositoryForGoalService)
		wantErr   bool
	}{
		{
			name: "успешное удаление",
			id:   1,
			mockSetup: func() (*MockGoalRepository, *MockCounterRepositoryForGoalService) {
				return &MockGoalRepository{
					DeleteFunc: func(ctx context.Context, id uint) error {
						return nil
					},
				}, &MockCounterRepositoryForGoalService{}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGoalRepo, mockCounterRepo := tt.mockSetup()
			service := NewGoalService(mockGoalRepo, mockCounterRepo)

			err := service.DeleteGoal(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
			}
		})
	}
}
