package services

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/suprt/planica_bi/backend/internal/middleware"
	"github.com/suprt/planica_bi/backend/internal/models"
)

// MockUserRepository implements UserRepositoryInterface for testing
type MockUserRepository struct {
	GetByEmailFunc      func(ctx context.Context, email string) (*models.User, error)
	GetByIDFunc         func(ctx context.Context, id uint) (*models.User, error)
	CreateFunc          func(ctx context.Context, user *models.User) error
	UpdateFunc          func(ctx context.Context, user *models.User) error
	DeleteFunc          func(ctx context.Context, userID uint) error
	UpdateLastLoginFunc func(ctx context.Context, userID uint) error
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, userID uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, userID)
	}
	return nil
}

// Stub methods for remaining UserRepositoryInterface methods
func (m *MockUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetAllPaginated(ctx context.Context, pagination *middleware.Pagination) ([]models.User, int64, error) {
	return nil, 0, nil
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID uint) error {
	if m.UpdateLastLoginFunc != nil {
		return m.UpdateLastLoginFunc(ctx, userID)
	}
	return nil
}

func (m *MockUserRepository) GetUserProjectRole(ctx context.Context, userID, projectID uint) (*models.UserProjectRole, error) {
	return nil, nil
}

func (m *MockUserRepository) GetUserProjects(ctx context.Context, userID uint) ([]models.UserProjectRole, error) {
	return nil, nil
}

func (m *MockUserRepository) GetProjectUsers(ctx context.Context, projectID uint) ([]models.UserProjectRole, error) {
	return nil, nil
}

func (m *MockUserRepository) AssignRole(ctx context.Context, userID, projectID uint, role string) error {
	return nil
}

func (m *MockUserRepository) UpdateRole(ctx context.Context, userID, projectID uint, role string) error {
	return nil
}

func (m *MockUserRepository) RemoveRole(ctx context.Context, userID, projectID uint) error {
	return nil
}

func (m *MockUserRepository) HasProjectAccess(ctx context.Context, userID, projectID uint) (bool, error) {
	return false, nil
}

func (m *MockUserRepository) IsAdmin(ctx context.Context, userID uint) (bool, error) {
	return false, nil
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name        string
		req         *RegisterRequest
		mockSetup   func() *MockUserRepository
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешная регистрация",
			req:  &RegisterRequest{Name: "Test User", Email: "test@example.com", Password: "password123"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{
					GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
						return nil, nil // User doesn't exist
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
			req:  &RegisterRequest{Name: "", Email: "test@example.com", Password: "password123"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{}
			},
			wantErr:     true,
			wantErrText: "name is required",
		},
		{
			name: "пустой email",
			req:  &RegisterRequest{Name: "Test", Email: "", Password: "password123"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{}
			},
			wantErr:     true,
			wantErrText: "email is required",
		},
		{
			name: "пустой пароль",
			req:  &RegisterRequest{Name: "Test", Email: "test@example.com", Password: ""},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{}
			},
			wantErr:     true,
			wantErrText: "password is required",
		},
		{
			name: "короткий пароль",
			req:  &RegisterRequest{Name: "Test", Email: "test@example.com", Password: "short"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{}
			},
			wantErr:     true,
			wantErrText: "password must be at least 8 characters",
		},
		{
			name: "пользователь уже существует",
			req:  &RegisterRequest{Name: "Test", Email: "test@example.com", Password: "password123"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{
					GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
						return &models.User{ID: 1, Email: "test@example.com"}, nil
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
			authService := NewAuthService(mockRepo, "test-secret-key-12345678901234567890123", time.Hour*24)

			resp, err := authService.Register(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
					return
				}
				if tt.wantErrText != "" && err.Error() != tt.wantErrText {
					t.Errorf("ожидалась ошибка '%s', но получили '%s'", tt.wantErrText, err.Error())
				}
				if resp != nil {
					t.Errorf("ожидался nil ответ, но получили %+v", resp)
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if resp == nil {
					t.Errorf("ожидался ответ, но получили nil")
				}
				if resp.Token == "" {
					t.Errorf("ожидался токен, но получили пустую строку")
				}
				if resp.User == nil {
					t.Errorf("ожидался пользователь, но получили nil")
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	// Generate a proper bcrypt hash for "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("не удалось создать хеш пароля: %v", err)
	}

	tests := []struct {
		name        string
		req         *LoginRequest
		mockSetup   func() *MockUserRepository
		wantErr     bool
		wantErrText string
	}{
		{
			name: "успешный вход",
			req:  &LoginRequest{Email: "test@example.com", Password: "password123"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{
					GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
						return &models.User{
							ID:       1,
							Email:    email,
							Name:     "Test User",
							Password: string(hashedPassword),
							IsActive: true,
						}, nil
					},
					UpdateLastLoginFunc: func(ctx context.Context, userID uint) error {
						return nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "пустой email",
			req:  &LoginRequest{Email: "", Password: "password123"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{}
			},
			wantErr:     true,
			wantErrText: "email is required",
		},
		{
			name: "пустой пароль",
			req:  &LoginRequest{Email: "test@example.com", Password: ""},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{}
			},
			wantErr:     true,
			wantErrText: "password is required",
		},
		{
			name: "пользователь не найден",
			req:  &LoginRequest{Email: "notfound@example.com", Password: "password123"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{
					GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
						return nil, nil
					},
				}
			},
			wantErr:     true,
			wantErrText: "invalid email or password",
		},
		{
			name: "неверный пароль",
			req:  &LoginRequest{Email: "test@example.com", Password: "wrongpassword"},
			mockSetup: func() *MockUserRepository {
				wrongHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				return &MockUserRepository{
					GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
						return &models.User{
							ID:       1,
							Email:    email,
							Password: string(wrongHash),
							IsActive: true,
						}, nil
					},
				}
			},
			wantErr:     true,
			wantErrText: "invalid email or password",
		},
		{
			name: "неактивный пользователь",
			req:  &LoginRequest{Email: "test@example.com", Password: "password123"},
			mockSetup: func() *MockUserRepository {
				return &MockUserRepository{
					GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
						return &models.User{
							ID:       1,
							Email:    email,
							Password: string(hashedPassword),
							IsActive: false,
						}, nil
					},
				}
			},
			wantErr:     true,
			wantErrText: "user account is inactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.mockSetup()
			authService := NewAuthService(mockRepo, "test-secret-key-12345678901234567890123", time.Hour*24)

			resp, err := authService.Login(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
					return
				}
				if tt.wantErrText != "" && err.Error() != tt.wantErrText {
					t.Errorf("ожидалась ошибка '%s', но получили '%s'", tt.wantErrText, err.Error())
				}
				if resp != nil {
					t.Errorf("ожидался nil ответ, но получили %+v", resp)
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if resp == nil {
					t.Errorf("ожидался ответ, но получили nil")
				}
				if resp.Token == "" {
					t.Errorf("ожидался токен, но получили пустую строку")
				}
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockRepo := &MockUserRepository{}
	authService := NewAuthService(mockRepo, "test-secret-key-12345678901234567890123", time.Hour*24)

	// First create a valid token
	userID := uint(123)
	token, err := authService.generateToken(userID)
	if err != nil {
		t.Fatalf("не удалось создать токен: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		wantID    uint
		wantErr   bool
		wantErrID uint
	}{
		{
			name:    "валидный токен",
			token:   token,
			wantID:  userID,
			wantErr: false,
		},
		{
			name:      "пустой токен",
			token:     "",
			wantErr:   true,
			wantErrID: 0,
		},
		{
			name:      "невалидный токен",
			token:     "invalid.token.here",
			wantErr:   true,
			wantErrID: 0,
		},
		{
			name:      "токен с неправильной подписью",
			token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImV4cCI6OTk5OTk5OTk5OX9.invalid_signature",
			wantErr:   true,
			wantErrID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := authService.ValidateToken(tt.token)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получили nil")
				}
				if gotID != tt.wantErrID {
					t.Errorf("ожидался user_id = %d, но получили %d", tt.wantErrID, gotID)
				}
			} else {
				if err != nil {
					t.Errorf("не ожидалась ошибка, но получили: %v", err)
				}
				if gotID != tt.wantID {
					t.Errorf("ожидался user_id = %d, но получили %d", tt.wantID, gotID)
				}
			}
		})
	}
}
