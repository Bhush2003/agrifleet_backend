package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MaintenanceRepository interface {
	Create(ctx context.Context, m *models.Maintenance) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Maintenance, error)
	List(ctx context.Context, machineID *uuid.UUID, offset, limit int) ([]models.Maintenance, int64, error)
	Update(ctx context.Context, m *models.Maintenance) error
	ListUpcoming(ctx context.Context, days int) ([]models.Maintenance, error)
	ListOverdue(ctx context.Context) ([]models.Maintenance, error)
}

type maintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) MaintenanceRepository {
	return &maintenanceRepository{db: db}
}

func (r *maintenanceRepository) Create(ctx context.Context, m *models.Maintenance) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("maintenanceRepository.Create: %w", err)
	}
	return nil
}

func (r *maintenanceRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Maintenance, error) {
	var m models.Maintenance
	if err := r.db.WithContext(ctx).Preload("Machine").First(&m, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("maintenanceRepository.FindByID: %w", err)
	}
	return &m, nil
}

func (r *maintenanceRepository) List(ctx context.Context, machineID *uuid.UUID, offset, limit int) ([]models.Maintenance, int64, error) {
	var records []models.Maintenance
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Maintenance{})
	if machineID != nil {
		q = q.Where("machine_id = ?", *machineID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("maintenanceRepository.List count: %w", err)
	}
	if err := q.Preload("Machine").Offset(offset).Limit(limit).Order("date DESC").Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("maintenanceRepository.List: %w", err)
	}
	return records, total, nil
}

func (r *maintenanceRepository) Update(ctx context.Context, m *models.Maintenance) error {
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("maintenanceRepository.Update: %w", err)
	}
	return nil
}

func (r *maintenanceRepository) ListUpcoming(ctx context.Context, days int) ([]models.Maintenance, error) {
	var records []models.Maintenance
	deadline := time.Now().AddDate(0, 0, days)
	err := r.db.WithContext(ctx).
		Preload("Machine").
		Where("next_due_date BETWEEN ? AND ? AND status != ?", time.Now(), deadline, models.MaintenanceCompleted).
		Order("next_due_date ASC").
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("maintenanceRepository.ListUpcoming: %w", err)
	}
	return records, nil
}

func (r *maintenanceRepository) ListOverdue(ctx context.Context) ([]models.Maintenance, error) {
	var records []models.Maintenance
	err := r.db.WithContext(ctx).
		Preload("Machine").
		Where("next_due_date < ? AND status != ?", time.Now(), models.MaintenanceCompleted).
		Order("next_due_date ASC").
		Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("maintenanceRepository.ListOverdue: %w", err)
	}
	return records, nil
}
