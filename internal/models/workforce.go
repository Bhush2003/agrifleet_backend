package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DriverStatus string
type WorkerStatus string
type AttendanceStatus string

const (
	DriverStatusActive   DriverStatus = "active"
	DriverStatusInactive DriverStatus = "inactive"
	DriverStatusOnLeave  DriverStatus = "on_leave"

	WorkerStatusActive   WorkerStatus = "active"
	WorkerStatusInactive WorkerStatus = "inactive"

	AttendancePresent  AttendanceStatus = "present"
	AttendanceAbsent   AttendanceStatus = "absent"
	AttendanceHalfDay  AttendanceStatus = "half_day"
	AttendanceOnLeave  AttendanceStatus = "on_leave"
)

// Driver links a User to a machine with license information.
type Driver struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	LicenseNumber string         `gorm:"type:varchar(50);uniqueIndex" json:"license_number"`
	LicenseExpiry time.Time      `json:"license_expiry"`
	MachineID     *uuid.UUID     `gorm:"type:uuid" json:"machine_id,omitempty"`
	Status        DriverStatus   `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	User    User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Machine *Machine `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
}

// Worker represents a field laborer (non-driver).
type Worker struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Role      string         `gorm:"type:varchar(50)" json:"role"`
	DailyWage float64        `gorm:"not null;default:0" json:"daily_wage"`
	JoinDate  time.Time      `json:"join_date"`
	Status    WorkerStatus   `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	Notes     string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Attendance records daily attendance for a worker.
type Attendance struct {
	ID         uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WorkerID   uuid.UUID        `gorm:"type:uuid;not null;index" json:"worker_id"`
	Date       time.Time        `gorm:"not null;index" json:"date"`
	CheckIn    *time.Time       `json:"check_in,omitempty"`
	CheckOut   *time.Time       `json:"check_out,omitempty"`
	Status     AttendanceStatus `gorm:"type:varchar(20);not null" json:"status"`
	WageEarned float64          `gorm:"default:0" json:"wage_earned"`
	Notes      string           `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`

	Worker Worker `gorm:"foreignKey:WorkerID" json:"worker,omitempty"`
}
