package models

import (
	"time"

	"github.com/google/uuid"
)

type ApprovalStatus string
type ApprovalType string

const (
	ApprovalPending    ApprovalStatus = "pending"
	ApprovalApproved   ApprovalStatus = "approved"
	ApprovalRejected   ApprovalStatus = "rejected"
	ApprovalReturned   ApprovalStatus = "returned"

	ApprovalAttendance  ApprovalType = "attendance"
	ApprovalSalary      ApprovalType = "salary"
	ApprovalDiesel      ApprovalType = "diesel"
	ApprovalExpense     ApprovalType = "expense"
	ApprovalMaintenance ApprovalType = "maintenance"
	ApprovalPayment     ApprovalType = "payment"
)

// ApprovalRequest is the centralized approval workflow entity.
type ApprovalRequest struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Type        ApprovalType   `gorm:"type:varchar(30);not null" json:"type"`
	Status      ApprovalStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	ReferenceID *uuid.UUID     `gorm:"type:uuid" json:"reference_id,omitempty"`
	ReferenceType string       `gorm:"type:varchar(50)" json:"reference_type,omitempty"`
	Amount      float64        `gorm:"default:0" json:"amount"`
	SubmittedBy uuid.UUID      `gorm:"type:uuid;not null" json:"submitted_by"`
	ReviewedBy  *uuid.UUID     `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewNote  string         `gorm:"type:text" json:"review_note,omitempty"`
	ReviewedAt  *time.Time     `json:"reviewed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	Submitter User  `gorm:"foreignKey:SubmittedBy" json:"submitter,omitempty"`
	Reviewer  *User `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}
