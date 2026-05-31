package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransportTrip records a single transportation trip.
type TransportTrip struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	JobID        *uuid.UUID     `gorm:"type:uuid;index" json:"job_id,omitempty"`
	TruckID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"truck_id"`
	DriverID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"driver_id"`
	Date         time.Time      `gorm:"not null;index" json:"date"`
	FromLocation string         `gorm:"type:varchar(200)" json:"from_location"`
	ToLocation   string         `gorm:"type:varchar(200)" json:"to_location"`
	Loads        int            `gorm:"default:1" json:"loads"`
	Tonnage      float64        `gorm:"default:0" json:"tonnage"`
	RatePerTon   float64        `gorm:"default:0" json:"rate_per_ton"`
	TotalAmount  float64        `gorm:"default:0" json:"total_amount"`
	Notes        string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Job    *HarvestingJob `gorm:"foreignKey:JobID" json:"job,omitempty"`
	Truck  Machine        `gorm:"foreignKey:TruckID" json:"truck,omitempty"`
	Driver Driver         `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
}
