package repositories

import (
	"context"
	"fmt"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	List(ctx context.Context, offset, limit int) ([]models.User, int64, error)
	SaveRefreshToken(ctx context.Context, rt *models.RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("userRepository.Create: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("userRepository.FindByID: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "phone = ?", phone).Error; err != nil {
		return nil, fmt.Errorf("userRepository.FindByPhone: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("userRepository.Update: %w", err)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("userRepository.List count: %w", err)
	}
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("userRepository.List: %w", err)
	}
	return users, total, nil
}

func (r *userRepository) SaveRefreshToken(ctx context.Context, rt *models.RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(rt).Error; err != nil {
		return fmt.Errorf("userRepository.SaveRefreshToken: %w", err)
	}
	return nil
}

func (r *userRepository) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	if err := r.db.WithContext(ctx).Preload("User").First(&rt, "token = ?", token).Error; err != nil {
		return nil, fmt.Errorf("userRepository.FindRefreshToken: %w", err)
	}
	return &rt, nil
}

func (r *userRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).Delete(&models.RefreshToken{}, "token = ?", token).Error; err != nil {
		return fmt.Errorf("userRepository.DeleteRefreshToken: %w", err)
	}
	return nil
}

func (r *userRepository) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.RefreshToken{}, "user_id = ?", userID).Error; err != nil {
		return fmt.Errorf("userRepository.DeleteUserRefreshTokens: %w", err)
	}
	return nil
}
