package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMode string
type ExpenseCategory string
type RevenueSource string
type ProjectStatus string

const (
	PaymentCash   PaymentMode = "cash"
	PaymentBank   PaymentMode = "bank_transfer"
	PaymentMobile PaymentMode = "mobile_money"
	PaymentCredit PaymentMode = "credit"

	ExpenseFuel        ExpenseCategory = "fuel"
	ExpenseMaintenance ExpenseCategory = "maintenance"
	ExpenseSalary      ExpenseCategory = "salary"
	ExpenseTransport   ExpenseCategory = "transport"
	ExpenseOther       ExpenseCategory = "other"

	RevenueHarvesting  RevenueSource = "harvesting"
	RevenueTransport   RevenueSource = "transport"
	RevenueProject     RevenueSource = "project"
	RevenueOther       RevenueSource = "other"

	ProjectActive    ProjectStatus = "active"
	ProjectCompleted ProjectStatus = "completed"
	ProjectCancelled ProjectStatus = "cancelled"
)

// Expense records any operational cost.
type Expense struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Category    ExpenseCategory `gorm:"type:varchar(30);not null" json:"category"`
	Amount      float64         `gorm:"not null" json:"amount"`
	Date        time.Time       `gorm:"not null;index" json:"date"`
	MachineID   *uuid.UUID      `gorm:"type:uuid" json:"machine_id,omitempty"`
	Description string          `gorm:"type:text" json:"description"`
	PaymentMode PaymentMode     `gorm:"type:varchar(20)" json:"payment_mode"`
	ReferenceID *uuid.UUID      `gorm:"type:uuid" json:"reference_id,omitempty"`
	CreatedBy   uuid.UUID       `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`

	Machine *Machine `gorm:"foreignKey:MachineID" json:"machine,omitempty"`
}

// Revenue records income from any source.
type Revenue struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SourceType  RevenueSource  `gorm:"type:varchar(30);not null" json:"source_type"`
	SourceID    *uuid.UUID     `gorm:"type:uuid" json:"source_id,omitempty"`
	Amount      float64        `gorm:"not null" json:"amount"`
	Date        time.Time      `gorm:"not null;index" json:"date"`
	Description string         `gorm:"type:text" json:"description"`
	PaymentMode PaymentMode    `gorm:"type:varchar(20)" json:"payment_mode"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Project represents a client contract/project.
type Project struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name       string         `gorm:"type:varchar(150);not null" json:"name"`
	ClientName string         `gorm:"type:varchar(150)" json:"client_name"`
	StartDate  time.Time      `json:"start_date"`
	EndDate    *time.Time     `json:"end_date,omitempty"`
	TotalValue float64        `gorm:"default:0" json:"total_value"`
	Status     ProjectStatus  `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	Notes      string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
