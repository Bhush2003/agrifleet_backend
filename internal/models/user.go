package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleOwner      UserRole = "owner"
	RoleManager    UserRole = "manager"
	RoleSupervisor UserRole = "supervisor"
	RoleDriver     UserRole = "driver"
	RoleAccountant UserRole = "accountant"
)

// User represents a system user with role-based access.
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"type:varchar(100);not null" json:"name"`
	Email        string         `gorm:"type:varchar(150);uniqueIndex" json:"email"`
	Phone        string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Role         UserRole       `gorm:"type:varchar(20);not null;default:'driver'" json:"role"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	AvatarURL    string         `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// RefreshToken stores JWT refresh tokens for session management.
type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(500);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
}
