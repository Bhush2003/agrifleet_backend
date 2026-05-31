package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockMovementType string

const (
	StockIn  StockMovementType = "in"
	StockOut StockMovementType = "out"
)

// SparePart represents an inventory item (spare part or consumable).
type SparePart struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string         `gorm:"type:varchar(150);not null" json:"name"`
	Category      string         `gorm:"type:varchar(80)" json:"category"`
	Unit          string         `gorm:"type:varchar(30)" json:"unit"`
	MinStockLevel float64        `gorm:"default:0" json:"min_stock_level"`
	CurrentStock  float64        `gorm:"default:0" json:"current_stock"`
	UnitCost      float64        `gorm:"default:0" json:"unit_cost"`
	Notes         string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// StockMovement records every stock-in or stock-out transaction.
type StockMovement struct {
	ID            uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SparePartID   uuid.UUID         `gorm:"type:uuid;not null;index" json:"spare_part_id"`
	Type          StockMovementType `gorm:"type:varchar(10);not null" json:"type"`
	Quantity      float64           `gorm:"not null" json:"quantity"`
	ReferenceID   *uuid.UUID        `gorm:"type:uuid" json:"reference_id,omitempty"`
	ReferenceType string            `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	Date          time.Time         `gorm:"not null" json:"date"`
	Notes         string            `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy     uuid.UUID         `gorm:"type:uuid" json:"created_by"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`

	SparePart SparePart `gorm:"foreignKey:SparePartID" json:"spare_part,omitempty"`
}
