package services

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/suprt/planica_bi/backend/internal/middleware"
	"github.com/suprt/planica_bi/backend/internal/models"
)

// MockUserRepositoryForUserService implements UserRepositoryInterface for user service testing
type MockUserRepositoryForUserService struct {
	CreateFunc           func(ctx context.Context, user *models.User) error
	GetByEmailFunc       func(ctx context.Context, email string) (*models.User, error)
	GetByIDFunc          func(ctx context.Context, id uint) (*models.User, error)
	GetAllFunc           func(ctx context.Context) ([]models.User, error)
	GetAllPaginatedFunc  func(ctx context.Context, pagination *middleware.Pagination) ([]models.User, int64, error)
	UpdateFunc           func(ctx context.Context, user *models.User) error
	DeleteFunc           func(ctx context.Context, userID uint) error
	UpdateLastLoginFunc  func(ctx context.Context, userID uint) error
	GetUserProjectsFunc  func(ctx context.Context, userID uint) ([]models.UserProjectRole, error)
	GetProjectUsersFunc  func(ctx context.Context, projectID uint) ([]models.UserProjectRole, error)
	AssignRoleFunc       func(ctx context.Context, userID, projectID uint, role string) error
	UpdateRoleFunc       func(ctx context.Context, userID, projectID uint, role string) error
	RemoveRoleFunc       func(ctx context.Context, userID, projectID uint) error
	HasProjectAccessFunc func(ctx context.Context, userID, projectID uint) (bool, error)
	IsAdminFunc          func(ctx context.Context, userID uint) (bool, error)
}

func (m *MockUserRepositoryForUserService) Create(ctx context.Context, user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepositoryForUserService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepositoryForUserService) GetAll(ctx context.Context) ([]models.User, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockUserRepositoryForUserService) GetAllPaginated(ctx context.Context, pagination *middleware.Pagination) ([]models.User, int64, error) {
	if m.GetAllPaginatedFunc != nil {
		return m.GetAllPaginatedFunc(ctx, pagination)
	}
	return nil, 0, nil
}

func (m *MockUserRepositoryForUserService) Update(ctx context.Context, user *models.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) Delete(ctx context.Context, userID uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, userID)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) UpdateLastLogin(ctx context.Context, userID uint) error {
	if m.UpdateLastLoginFunc != nil {
		return m.UpdateLastLoginFunc(ctx, userID)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) GetUserProjects(ctx context.Context, userID uint) ([]models.UserProjectRole, error) {
	if m.GetUserProjectsFunc != nil {
		return m.GetUserProjectsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockUserRepositoryForUserService) GetProjectUsers(ctx context.Context, projectID uint) ([]models.UserProjectRole, error) {
	if m.GetProjectUsersFunc != nil {
		return m.GetProjectUsersFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockUserRepositoryForUserService) AssignRole(ctx context.Context, userID, projectID uint, role string) error {
	if m.AssignRoleFunc != nil {
		return m.AssignRoleFunc(ctx, userID, projectID, role)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) UpdateRole(ctx context.Context, userID, projectID uint, role string) error {
	if m.UpdateRoleFunc != nil {
		return m.UpdateRoleFunc(ctx, userID, projectID, role)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) RemoveRole(ctx context.Context, userID, projectID uint) error {
	if m.RemoveRoleFunc != nil {
		return m.RemoveRoleFunc(ctx, userID, projectID)
	}
	return nil
}

func (m *MockUserRepositoryForUserService) HasProjectAccess(ctx context.Context, userID, projectID uint) (bool, error) {
	if m.HasProjectAccessFunc != nil {
		return m.HasProjectAccessFunc(ctx, userID, projectID)
	}
	return false, nil
}

func (m *MockUserRepositoryForUserService) IsAdmin(ctx context.Context, userID uint) (bool, error) {
	if m.IsAdminFunc != nil {
		return m.IsAdminFunc(ctx, userID)
	}
	return false, nil
}

func (m *MockUserRepositoryForUserService) GetUserProjectRole(ctx context.Context, userID, projectID uint) (*models.UserProjectRole, error) {
	return nil, nil
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		req         *CreateUserRequest
		mockSetup   func() *MockUserRepositoryForUserService
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешное создание пользователя",
			req: &CreateUserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				IsActive: true,
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
						return nil, nil
					},
					CreateFunc: func(ctx context.Context, user *models.User) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "пустое имя",
			req: &CreateUserRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{}
			},
			wantErr:     true,
			wantErrText: "name is required",
		},
		{
			name: "пустой email",
			req: &CreateUserRequest{
				Name:     "Test",
				Email:    "",
				Password: "password123",
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{}
			},
			wantErr:     true,
			wantErrText: "email is required",
		},
		{
			name: "пустой пароль",
			req: &CreateUserRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Password: "",
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{}
			},
			wantErr:     true,
			wantErrText: "password is required",
		},
		{
			name: "короткий пароль",
			req: &CreateUserRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Password: "short",
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{}
			},
			wantErr:     true,
			wantErrText: "password must be at least 8 characters",
		},
		{
			name: "пользователь уже существует",
			req: &CreateUserRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
						return &models.User{ID: 1, Email: email}, nil
					},
				}
			},
			wantErr:     true,
			wantErrText: "user with this email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewUserService(mockRepo)

			resp, err := service.CreateUser(context.Background(), tt.req)

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
				if resp == nil {
					t.Errorf("ожидался ответ, но получили nil")
				}
				if resp.Email != tt.req.Email {
					t.Errorf("ожидался email = %s, но получили %s", tt.req.Email, resp.Email)
				}
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		userID      uint
		req         *UpdateUserRequest
		mockSetup   func() *MockUserRepositoryForUserService
		wantErr     bool
		wantErrText string
	}{
		{
			name:   "успешное обновление",
			userID: 1,
			req: &UpdateUserRequest{
				Name: stringPtr("Updated Name"),
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.User, error) {
						return &models.User{ID: id, Name: "Old Name", Password: string(hashedPassword)}, nil
					},
					UpdateFunc: func(ctx context.Context, user *models.User) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "пользователь не найден",
			userID: 999,
			req: &UpdateUserRequest{
				Name: stringPtr("Updated Name"),
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.User, error) {
						return nil, nil
					},
				}
			},
			wantErr: true,
		},
		{
			name:   "пустое имя при обновлении",
			userID: 1,
			req: &UpdateUserRequest{
				Name: stringPtr(""),
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.User, error) {
						return &models.User{ID: id, Name: "Old Name", Password: string(hashedPassword)}, nil
					},
				}
			},
			wantErr:     true,
			wantErrText: "name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewUserService(mockRepo)

			resp, err := service.UpdateUser(context.Background(), tt.userID, tt.req)

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
				if resp == nil {
					t.Errorf("ожидался ответ, но получили nil")
				}
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		mockSetup func() *MockUserRepositoryForUserService
		wantErr   bool
	}{
		{
			name:   "успешное удаление",
			userID: 1,
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.User, error) {
						return &models.User{ID: id}, nil
					},
					DeleteFunc: func(ctx context.Context, userID uint) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:   "пользователь не найден",
			userID: 999,
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.User, error) {
						return nil, nil
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewUserService(mockRepo)

			err := service.DeleteUser(context.Background(), tt.userID)

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

func TestUserService_AssignRole(t *testing.T) {
	tests := []struct {
		name        string
		req         *AssignRoleRequest
		mockSetup   func() *MockUserRepositoryForUserService
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешное назначение роли",
			req: &AssignRoleRequest{
				UserID:    1,
				ProjectID: 2,
				Role:      "client",
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.User, error) {
						return &models.User{ID: id}, nil
					},
					AssignRoleFunc: func(ctx context.Context, userID, projectID uint, role string) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "невалидная роль",
			req: &AssignRoleRequest{
				UserID:    1,
				ProjectID: 2,
				Role:      "invalid",
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{}
			},
			wantErr:     true,
			wantErrText: "invalid role, must be 'admin', 'manager', or 'client'",
		},
		{
			name: "пользователь не найден",
			req: &AssignRoleRequest{
				UserID:    999,
				ProjectID: 2,
				Role:      "client",
			},
			mockSetup: func() *MockUserRepositoryForUserService {
				return &MockUserRepositoryForUserService{
					GetByIDFunc: func(ctx context.Context, id uint) (*models.User, error) {
						return nil, nil
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			service := NewUserService(mockRepo)

			err := service.AssignRole(context.Background(), tt.req)

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

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
