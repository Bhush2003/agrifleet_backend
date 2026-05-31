package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobStatus string

const (
	JobUpcoming  JobStatus = "upcoming"
	JobActive    JobStatus = "active"
	JobCompleted JobStatus = "completed"
	JobCancelled JobStatus = "cancelled"
)

// HarvestingJob represents a harvesting contract/operation.
type HarvestingJob struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FarmName   string         `gorm:"type:varchar(150);not null" json:"farm_name"`
	Location   string         `gorm:"type:varchar(200)" json:"location"`
	CropType   string         `gorm:"type:varchar(80)" json:"crop_type"`
	StartDate  time.Time      `json:"start_date"`
	EndDate    *time.Time     `json:"end_date,omitempty"`
	TotalAcres float64        `gorm:"default:0" json:"total_acres"`
	Status     JobStatus      `gorm:"type:varchar(20);not null;default:'upcoming'" json:"status"`
	MachineID  *uuid.UUID     `gorm:"type:uuid" json:"machine_id,omitempty"`
	ProjectID  *uuid.UUID     `gorm:"type:uuid" json:"project_id,omitempty"`
	Notes      string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Machine *Machine `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// HarvestingLog records daily progress for a harvesting job.
type HarvestingLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	JobID     uuid.UUID `gorm:"type:uuid;not null;index" json:"job_id"`
	Date      time.Time `gorm:"not null" json:"date"`
	MachineID uuid.UUID `gorm:"type:uuid;not null" json:"machine_id"`
	DriverID  uuid.UUID `gorm:"type:uuid" json:"driver_id"`
	AcresDone float64   `gorm:"default:0" json:"acres_done"`
	Tonnage   float64   `gorm:"default:0" json:"tonnage"`
	Notes     string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Job     HarvestingJob `gorm:"foreignKey:JobID" json:"job,omitempty"`
	Machine Machine       `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
	Driver  Driver        `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
}
