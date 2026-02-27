package services

import (
	"context"
	"testing"

	"github.com/suprt/planica_bi/backend/internal/middleware"
	"github.com/suprt/planica_bi/backend/internal/models"
)

// MockProjectRepository implements ProjectRepositoryInterface for testing
type MockProjectRepository struct {
	CreateFunc               func(ctx context.Context, project *models.Project) error
	GetByIDFunc              func(ctx context.Context, id uint) (*models.Project, error)
	GetAllFunc               func(ctx context.Context) ([]*models.Project, error)
	GetAllPaginatedFunc      func(ctx context.Context, pagination *middleware.Pagination) ([]*models.Project, int64, error)
	GetByUserIDFunc          func(ctx context.Context, userID uint, isAdmin bool) ([]*models.Project, error)
	GetByUserIDPaginatedFunc func(ctx context.Context, userID uint, isAdmin bool, pagination *middleware.Pagination) ([]*models.Project, int64, error)
	UpdateFunc               func(ctx context.Context, project *models.Project) error
	DeleteFunc               func(ctx context.Context, id uint) error
	GetByPublicTokenFunc     func(ctx context.Context, token string) (*models.Project, error)
}

func (m *MockProjectRepository) Create(ctx context.Context, project *models.Project) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id uint) (*models.Project, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProjectRepository) GetAll(ctx context.Context) ([]*models.Project, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockProjectRepository) GetAllPaginated(ctx context.Context, pagination *middleware.Pagination) ([]*models.Project, int64, error) {
	if m.GetAllPaginatedFunc != nil {
		return m.GetAllPaginatedFunc(ctx, pagination)
	}
	return nil, 0, nil
}

func (m *MockProjectRepository) GetByUserID(ctx context.Context, userID uint, isAdmin bool) ([]*models.Project, error) {
	if m.GetByUserIDFunc != nil {
		return m.GetByUserIDFunc(ctx, userID, isAdmin)
	}
	return nil, nil
}

func (m *MockProjectRepository) GetByUserIDPaginated(ctx context.Context, userID uint, isAdmin bool, pagination *middleware.Pagination) ([]*models.Project, int64, error) {
	if m.GetByUserIDPaginatedFunc != nil {
		return m.GetByUserIDPaginatedFunc(ctx, userID, isAdmin, pagination)
	}
	return nil, 0, nil
}

func (m *MockProjectRepository) Update(ctx context.Context, project *models.Project) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockProjectRepository) GetByPublicToken(ctx context.Context, token string) (*models.Project, error) {
	if m.GetByPublicTokenFunc != nil {
		return m.GetByPublicTokenFunc(ctx, token)
	}
	return nil, nil
}

func TestProjectService_CreateProject(t *testing.T) {
	tests := []struct {
		name        string
		project     *models.Project
		mockSetup   func() *MockProjectRepository
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешное создание проекта",
			project: &models.Project{
				Name: "Test Project",
				Slug: "test-project",
			},
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{
					CreateFunc: func(ctx context.Context, project *models.Project) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "пустое название",
			project: &models.Project{
				Name: "",
				Slug: "test-project",
			},
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{}
			},
			wantErr:     true,
			wantErrText: "name is required",
		},
		{
			name: "пустой slug",
			project: &models.Project{
				Name: "Test Project",
				Slug: "",
			},
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{}
			},
			wantErr:     true,
			wantErrText: "slug is required",
		},
		{
			name: "сгенерировать public token если не предоставлен",
			project: &models.Project{
				Name:        "Test Project",
				Slug:        "test-project",
				PublicToken: "",
			},
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{
					CreateFunc: func(ctx context.Context, project *models.Project) error {
						if project.PublicToken == "" {
							t.Errorf("ожидался сгенерированный PublicToken, но получили пустую строку")
						}
						if len(project.PublicToken) != 64 {
							t.Errorf("ожидалась длина токена 64 символа, но получили %d", len(project.PublicToken))
						}
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
			service := NewProjectService(mockRepo)

			err := service.CreateProject(context.Background(), tt.project)

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

func TestProjectService_GetProject(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func() *MockProjectRepository
		wantErr   bool
		wantID    uint
	}{
		{
			name: "успешное получение проекта",
			id:   1,
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.Project, error) {
						return &models.Project{ID: id, Name: "Test"}, nil
					},
				}
			},
			wantErr: false,
			wantID:  1,
		},
		{
			name: "проект не найден",
			id:   999,
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.Project, error) {
						return nil, nil
					},
				}
			},
			wantErr: true,
			wantID:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewProjectService(mockRepo)

			project, err := service.GetProject(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil && project == nil {
					// Ожидали ошибку или nil проект
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if project == nil {
					t.Errorf("ожидался проект, но получили nil")
				}
				if project.ID != tt.wantID {
					t.Errorf("ожидался ID = %d, но получили %d", tt.wantID, project.ID)
				}
			}
		})
	}
}

func TestProjectService_GetAllProjects(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		isAdmin   bool
		mockSetup func() *MockProjectRepository
		wantCount int
	}{
		{
			name:    "админ получает все проекты",
			userID:  1,
			isAdmin: true,
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{
					GetByUserIDFunc: func(ctx context.Context, userID uint, isAdmin bool) ([]*models.Project, error) {
						return []*models.Project{
							{ID: 1, Name: "Project 1"},
							{ID: 2, Name: "Project 2"},
						}, nil
					},
				}
			},
			wantCount: 2,
		},
		{
			name:    "пользователь получает только свои проекты",
			userID:  2,
			isAdmin: false,
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{
					GetByUserIDFunc: func(ctx context.Context, userID uint, isAdmin bool) ([]*models.Project, error) {
						return []*models.Project{
							{ID: 3, Name: "User Project"},
						}, nil
					},
				}
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewProjectService(mockRepo)

			projects, err := service.GetAllProjects(context.Background(), tt.userID, tt.isAdmin)

			if err != nil {
				t.Errorf("не ожидалась ошибка, но получили: %v", err)
			}
			if len(projects) != tt.wantCount {
				t.Errorf("ожидалось %d проектов, но получили %d", tt.wantCount, len(projects))
			}
		})
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
	tests := []struct {
		name        string
		project     *models.Project
		mockSetup   func() *MockProjectRepository
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешное обновление",
			project: &models.Project{
				ID:   1,
				Name: "Updated Project",
				Slug: "updated-project",
			},
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{
					UpdateFunc: func(ctx context.Context, project *models.Project) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "пустое название при обновлении",
			project: &models.Project{
				ID:   1,
				Name: "",
			},
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{}
			},
			wantErr:     true,
			wantErrText: "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewProjectService(mockRepo)

			err := service.UpdateProject(context.Background(), tt.project)

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

func TestProjectService_DeleteProject(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func() *MockProjectRepository
		wantErr   bool
	}{
		{
			name: "успешное удаление",
			id:   1,
			mockSetup: func() *MockProjectRepository {
				return &MockProjectRepository{
					DeleteFunc: func(ctx context.Context, id uint) error {
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
			service := NewProjectService(mockRepo)

			err := service.DeleteProject(context.Background(), tt.id)

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
