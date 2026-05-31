package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FuelEntry records a diesel/fuel fill-up for a machine.
type FuelEntry struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MachineID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"machine_id"`
	Date         time.Time      `gorm:"not null;index" json:"date"`
	Liters       float64        `gorm:"not null" json:"liters"`
	CostPerLiter float64        `gorm:"not null" json:"cost_per_liter"`
	TotalCost    float64        `gorm:"not null" json:"total_cost"`
	FilledBy     uuid.UUID      `gorm:"type:uuid" json:"filled_by"`
	Odometer     float64        `json:"odometer"`
	ReceiptURL   string         `gorm:"type:varchar(500)" json:"receipt_url,omitempty"`
	Notes        string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Machine  Machine `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
	FilledByUser User `gorm:"foreignKey:FilledBy" json:"filled_by_user,omitempty"`
}
