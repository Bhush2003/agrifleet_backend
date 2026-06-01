package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Farmer represents a client farmer who provides crops for harvesting.
type Farmer struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string         `gorm:"type:varchar(150);not null" json:"name"`
	Mobile          string         `gorm:"type:varchar(20)" json:"mobile"`
	Village         string         `gorm:"type:varchar(100)" json:"village"`
	ContractRate    float64        `gorm:"default:0" json:"contract_rate"`
	TotalTon        float64        `gorm:"default:0" json:"total_ton"`
	TotalPayment    float64        `gorm:"default:0" json:"total_payment"`
	PendingBalance  float64        `gorm:"default:0" json:"pending_balance"`
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedBy       uuid.UUID      `gorm:"type:uuid" json:"created_by"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedBy       *uuid.UUID     `gorm:"type:uuid" json:"deleted_by,omitempty"`
}

// FarmerLedger tracks all financial transactions for a farmer.
type FarmerLedger struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FarmerID    uuid.UUID `gorm:"type:uuid;not null;index" json:"farmer_id"`
	Type        string    `gorm:"type:varchar(30);not null" json:"type"` // ton_entry, payment, adjustment
	Amount      float64   `gorm:"not null" json:"amount"`
	TonQty      float64   `gorm:"default:0" json:"ton_qty"`
	Description string    `gorm:"type:text" json:"description"`
	ReferenceID *uuid.UUID `gorm:"type:uuid" json:"reference_id,omitempty"`
	Date        time.Time `gorm:"not null" json:"date"`
	CreatedBy   uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`

	Farmer Farmer `gorm:"foreignKey:FarmerID" json:"farmer,omitempty"`
}
