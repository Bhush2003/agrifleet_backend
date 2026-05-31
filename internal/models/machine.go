package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MachineStatus string
type MachineType string

const (
	MachineStatusActive      MachineStatus = "active"
	MachineStatusIdle        MachineStatus = "idle"
	MachineStatusMaintenance MachineStatus = "maintenance"
	MachineStatusRetired     MachineStatus = "retired"

	MachineTypeHarvester MachineType = "harvester"
	MachineTypeTractor   MachineType = "tractor"
	MachineTypeTruck     MachineType = "truck"
	MachineTypeLoader    MachineType = "loader"
	MachineTypeOther     MachineType = "other"
)

// Machine represents an agricultural machine in the fleet.
type Machine struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name            string         `gorm:"type:varchar(100);not null" json:"name"`
	Type            MachineType    `gorm:"type:varchar(30);not null" json:"type"`
	Model           string         `gorm:"type:varchar(100)" json:"model"`
	Year            int            `json:"year"`
	RegNumber       string         `gorm:"type:varchar(50);uniqueIndex" json:"reg_number"`
	Status          MachineStatus  `gorm:"type:varchar(20);not null;default:'idle'" json:"status"`
	CurrentLocation string         `gorm:"type:varchar(200)" json:"current_location"`
	ImageURL        string         `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	OwnerID         uuid.UUID      `gorm:"type:uuid;not null" json:"owner_id"`
	Notes           string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Owner User `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
}

// MachineLog records daily working hours and location for a machine.
type MachineLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MachineID    uuid.UUID `gorm:"type:uuid;not null;index" json:"machine_id"`
	Date         time.Time `gorm:"not null" json:"date"`
	WorkingHours float64   `gorm:"not null;default:0" json:"working_hours"`
	Location     string    `gorm:"type:varchar(200)" json:"location"`
	OperatorID   uuid.UUID `gorm:"type:uuid" json:"operator_id"`
	Notes        string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Machine  Machine `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
	Operator User    `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
}
