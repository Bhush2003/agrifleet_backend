package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Factory represents a processing factory that receives harvested crops.
type Factory struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string         `gorm:"type:varchar(150);not null" json:"name"`
	RatePerTon   float64        `gorm:"default:0" json:"rate_per_ton"`
	ContactName  string         `gorm:"type:varchar(100)" json:"contact_name"`
	ContactPhone string         `gorm:"type:varchar(20)" json:"contact_phone"`
	Address      string         `gorm:"type:varchar(300)" json:"address"`
	IsLocked     bool           `gorm:"default:false" json:"is_locked"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedBy    uuid.UUID      `gorm:"type:uuid" json:"created_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// WeightSlip records a weight measurement at the factory.
type WeightSlip struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FactoryID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"factory_id"`
	FarmerID   *uuid.UUID `gorm:"type:uuid" json:"farmer_id,omitempty"`
	MachineID  *uuid.UUID `gorm:"type:uuid" json:"machine_id,omitempty"`
	DriverID   *uuid.UUID `gorm:"type:uuid" json:"driver_id,omitempty"`
	SlipNumber string     `gorm:"type:varchar(50)" json:"slip_number"`
	GrossWeight float64   `gorm:"default:0" json:"gross_weight"`
	TareWeight  float64   `gorm:"default:0" json:"tare_weight"`
	NetWeight   float64   `gorm:"default:0" json:"net_weight"`
	RatePerTon  float64   `gorm:"default:0" json:"rate_per_ton"`
	TotalAmount float64   `gorm:"default:0" json:"total_amount"`
	SlipPhotoURL string   `gorm:"type:varchar(500)" json:"slip_photo_url,omitempty"`
	Date        time.Time `gorm:"not null" json:"date"`
	CreatedBy   uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Factory Factory `gorm:"foreignKey:FactoryID" json:"factory,omitempty"`
}
