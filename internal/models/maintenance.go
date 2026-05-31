package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaintenanceType   string
type MaintenanceStatus string

const (
	MaintenanceScheduled  MaintenanceType = "scheduled"
	MaintenanceBreakdown  MaintenanceType = "breakdown"
	MaintenancePreventive MaintenanceType = "preventive"

	MaintenancePending    MaintenanceStatus = "pending"
	MaintenanceInProgress MaintenanceStatus = "in_progress"
	MaintenanceCompleted  MaintenanceStatus = "completed"
	MaintenanceOverdue    MaintenanceStatus = "overdue"
)

// Maintenance records a maintenance event for a machine.
type Maintenance struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MachineID   uuid.UUID         `gorm:"type:uuid;not null;index" json:"machine_id"`
	Type        MaintenanceType   `gorm:"type:varchar(20);not null" json:"type"`
	Date        time.Time         `gorm:"not null" json:"date"`
	Description string            `gorm:"type:text;not null" json:"description"`
	Cost        float64           `gorm:"default:0" json:"cost"`
	NextDueDate *time.Time        `json:"next_due_date,omitempty"`
	Status      MaintenanceStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Technician  string            `gorm:"type:varchar(100)" json:"technician"`
	Notes       string            `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`

	Machine Machine `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
}
