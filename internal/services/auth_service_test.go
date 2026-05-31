package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/agrifleet/backend/internal/config"
	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ── Mock UserRepository ────────────────────────────────────────────────────────

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepo) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepo) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepo) List(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserRepo) SaveRefreshToken(ctx context.Context, rt *models.RefreshToken) error {
	args := m.Called(ctx, rt)
	return args.Error(0)
}

func (m *mockUserRepo) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *mockUserRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockUserRepo) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// ── Tests ──────────────────────────────────────────────────────────────────────

func testConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret-key",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           uuid.New(),
		Phone:        "+1234567890",
		PasswordHash: string(hash),
		Role:         models.RoleOwner,
		IsActive:     true,
	}

	repo := new(mockUserRepo)
	repo.On("FindByPhone", mock.Anything, "+1234567890").Return(user, nil)
	repo.On("SaveRefreshToken", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	svc := services.NewAuthService(repo, testConfig())
	tokens, err := svc.Login(context.Background(), services.LoginRequest{
		Phone:    "+1234567890",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Equal(t, user.ID, tokens.User.ID)
	repo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           uuid.New(),
		Phone:        "+1234567890",
		PasswordHash: string(hash),
		IsActive:     true,
	}

	repo := new(mockUserRepo)
	repo.On("FindByPhone", mock.Anything, "+1234567890").Return(user, nil)

	svc := services.NewAuthService(repo, testConfig())
	_, err := svc.Login(context.Background(), services.LoginRequest{
		Phone:    "+1234567890",
		Password: "wrongpassword",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindByPhone", mock.Anything, "+9999999999").Return(nil, gorm.ErrRecordNotFound)

	svc := services.NewAuthService(repo, testConfig())
	_, err := svc.Login(context.Background(), services.LoginRequest{
		Phone:    "+9999999999",
		Password: "anypassword",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:           uuid.New(),
		Phone:        "+1234567890",
		PasswordHash: string(hash),
		IsActive:     false,
	}

	repo := new(mockUserRepo)
	repo.On("FindByPhone", mock.Anything, "+1234567890").Return(user, nil)

	svc := services.NewAuthService(repo, testConfig())
	_, err := svc.Login(context.Background(), services.LoginRequest{
		Phone:    "+1234567890",
		Password: "password123",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deactivated")
}

func TestAuthService_Logout(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("DeleteRefreshToken", mock.Anything, "some-refresh-token").Return(nil)

	svc := services.NewAuthService(repo, testConfig())
	err := svc.Logout(context.Background(), "some-refresh-token")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAuthService_RefreshTokens_Expired(t *testing.T) {
	rt := &models.RefreshToken{
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // already expired
		User:      models.User{ID: uuid.New()},
	}

	repo := new(mockUserRepo)
	repo.On("FindRefreshToken", mock.Anything, "expired-token").Return(rt, nil)
	repo.On("DeleteRefreshToken", mock.Anything, "expired-token").Return(nil)

	svc := services.NewAuthService(repo, testConfig())
	_, err := svc.RefreshTokens(context.Background(), "expired-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestAuthService_RefreshTokens_InvalidToken(t *testing.T) {
	repo := new(mockUserRepo)
	repo.On("FindRefreshToken", mock.Anything, "bad-token").Return(nil, errors.New("not found"))

	svc := services.NewAuthService(repo, testConfig())
	_, err := svc.RefreshTokens(context.Background(), "bad-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token")
}
