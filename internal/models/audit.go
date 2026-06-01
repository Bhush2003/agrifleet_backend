package models

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog records every significant action in the system.
type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Action     string    `gorm:"type:varchar(100);not null" json:"action"`
	Module     string    `gorm:"type:varchar(50);not null" json:"module"`
	EntityID   *uuid.UUID `gorm:"type:uuid" json:"entity_id,omitempty"`
	EntityType string    `gorm:"type:varchar(50)" json:"entity_type,omitempty"`
	OldValue   string    `gorm:"type:jsonb" json:"old_value,omitempty"`
	NewValue   string    `gorm:"type:jsonb" json:"new_value,omitempty"`
	IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  string    `gorm:"type:varchar(300)" json:"user_agent,omitempty"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// MonthLock prevents editing of records for a locked month.
type MonthLock struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Month      int       `gorm:"not null" json:"month"`
	Year       int       `gorm:"not null" json:"year"`
	Module     string    `gorm:"type:varchar(50);not null" json:"module"` // attendance, salary, diesel, all
	IsLocked   bool      `gorm:"default:false" json:"is_locked"`
	LockedBy   *uuid.UUID `gorm:"type:uuid" json:"locked_by,omitempty"`
	LockedAt   *time.Time `json:"locked_at,omitempty"`
	UnlockedBy *uuid.UUID `gorm:"type:uuid" json:"unlocked_by,omitempty"`
	UnlockedAt *time.Time `json:"unlocked_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// RecycleBin stores soft-deleted records.
type RecycleBin struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EntityType  string    `gorm:"type:varchar(50);not null" json:"entity_type"`
	EntityID    uuid.UUID `gorm:"type:uuid;not null" json:"entity_id"`
	EntityData  string    `gorm:"type:jsonb;not null" json:"entity_data"`
	DeletedBy   uuid.UUID `gorm:"type:uuid;not null" json:"deleted_by"`
	DeletedAt   time.Time `gorm:"not null" json:"deleted_at"`
	RestoredBy  *uuid.UUID `gorm:"type:uuid" json:"restored_by,omitempty"`
	RestoredAt  *time.Time `json:"restored_at,omitempty"`
	IsRestored  bool      `gorm:"default:false" json:"is_restored"`

	Deleter  User  `gorm:"foreignKey:DeletedBy" json:"deleter,omitempty"`
}

// Notification stores in-app notifications.
type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title     string    `gorm:"type:varchar(200);not null" json:"title"`
	Body      string    `gorm:"type:text" json:"body"`
	Type      string    `gorm:"type:varchar(50)" json:"type"` // alert, approval, salary, diesel
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	ActionURL string    `gorm:"type:varchar(300)" json:"action_url,omitempty"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Task represents an assignable task for managers/drivers.
type Task struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      string     `gorm:"type:varchar(20);not null;default:'open'" json:"status"` // open, in_progress, completed, overdue
	Priority    string     `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	AssignedTo  *uuid.UUID `gorm:"type:uuid" json:"assigned_to,omitempty"`
	CreatedBy   uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Assignee *User `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	Creator  User  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}
