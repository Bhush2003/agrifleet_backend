package models

import (
	"time"

	"github.com/google/uuid"
)

type DocumentType string

const (
	DocAadhaar        DocumentType = "aadhaar"
	DocDrivingLicense DocumentType = "driving_license"
	DocRC             DocumentType = "rc"
	DocInsurance      DocumentType = "insurance"
	DocFitness        DocumentType = "fitness"
	DocPUC            DocumentType = "puc"
	DocOther          DocumentType = "other"
)

// Document tracks expiry-sensitive documents for employees and vehicles.
type Document struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EntityType  string       `gorm:"type:varchar(30);not null" json:"entity_type"` // employee, vehicle
	EntityID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"entity_id"`
	Type        DocumentType `gorm:"type:varchar(30);not null" json:"type"`
	DocumentNo  string       `gorm:"type:varchar(100)" json:"document_no"`
	FileURL     string       `gorm:"type:varchar(500)" json:"file_url,omitempty"`
	IssueDate   *time.Time   `json:"issue_date,omitempty"`
	ExpiryDate  *time.Time   `json:"expiry_date,omitempty"`
	IsExpired   bool         `gorm:"default:false" json:"is_expired"`
	AlertSent30 bool         `gorm:"default:false" json:"alert_sent_30"`
	AlertSent15 bool         `gorm:"default:false" json:"alert_sent_15"`
	AlertSent7  bool         `gorm:"default:false" json:"alert_sent_7"`
	UploadedBy  uuid.UUID    `gorm:"type:uuid" json:"uploaded_by"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Ledger is a unified ledger entry for any entity (employee, vehicle, farmer).
type Ledger struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EntityType  string    `gorm:"type:varchar(30);not null;index" json:"entity_type"` // employee, vehicle, farmer
	EntityID    uuid.UUID `gorm:"type:uuid;not null;index" json:"entity_id"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type"` // salary, advance, bonus, deduction, diesel, repair, income, payment
	Debit       float64   `gorm:"default:0" json:"debit"`
	Credit      float64   `gorm:"default:0" json:"credit"`
	Balance     float64   `gorm:"default:0" json:"balance"`
	Description string    `gorm:"type:text" json:"description"`
	ReferenceID *uuid.UUID `gorm:"type:uuid" json:"reference_id,omitempty"`
	Date        time.Time `gorm:"not null;index" json:"date"`
	CreatedBy   uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// SystemSetting stores key-value configuration.
type SystemSetting struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Key       string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Group     string    `gorm:"type:varchar(50)" json:"group"` // salary, company, whatsapp, backup
	UpdatedBy uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`
}
