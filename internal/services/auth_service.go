package services

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/config"
	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *models.User `json:"user"`
}

type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error)
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*TokenPair, error) {
	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("authService.Login: %w", err)
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return s.generateTokenPair(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.userRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return fmt.Errorf("authService.Logout: %w", err)
	}
	return nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	rt, err := s.userRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = s.userRepo.DeleteRefreshToken(ctx, refreshToken)
		return nil, fmt.Errorf("refresh token expired")
	}

	// Rotate: delete old, issue new
	_ = s.userRepo.DeleteRefreshToken(ctx, refreshToken)
	return s.generateTokenPair(ctx, &rt.User)
}

func (s *authService) generateTokenPair(ctx context.Context, user *models.User) (*TokenPair, error) {
	accessToken, err := pkgjwt.GenerateAccessToken(user.ID, string(user.Role), s.cfg.JWT.Secret, s.cfg.JWT.AccessExpiry)
	if err != nil {
		return nil, fmt.Errorf("authService.generateTokenPair access: %w", err)
	}

	rawRefresh := pkgjwt.GenerateRefreshToken()
	expiresAt := time.Now().Add(s.cfg.JWT.RefreshExpiry)

	rt := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     rawRefresh,
		ExpiresAt: expiresAt,
	}
	if err := s.userRepo.SaveRefreshToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("authService.generateTokenPair refresh: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
		ExpiresAt:    time.Now().Add(s.cfg.JWT.AccessExpiry),
		User:         user,
	}, nil
}
