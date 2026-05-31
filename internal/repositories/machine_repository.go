package repositories

import (
	"context"
	"fmt"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MachineRepository interface {
	Create(ctx context.Context, m *models.Machine) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Machine, error)
	List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Machine, int64, error)
	Update(ctx context.Context, m *models.Machine) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateLog(ctx context.Context, log *models.MachineLog) error
	ListLogs(ctx context.Context, machineID uuid.UUID, offset, limit int) ([]models.MachineLog, int64, error)
}

type machineRepository struct {
	db *gorm.DB
}

func NewMachineRepository(db *gorm.DB) MachineRepository {
	return &machineRepository{db: db}
}

func (r *machineRepository) Create(ctx context.Context, m *models.Machine) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("machineRepository.Create: %w", err)
	}
	return nil
}

func (r *machineRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Machine, error) {
	var m models.Machine
	if err := r.db.WithContext(ctx).Preload("Owner").First(&m, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("machineRepository.FindByID: %w", err)
	}
	return &m, nil
}

func (r *machineRepository) List(ctx context.Context, filters map[string]interface{}, offset, limit int) ([]models.Machine, int64, error) {
	var machines []models.Machine
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Machine{})
	for k, v := range filters {
		query = query.Where(k+" = ?", v)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("machineRepository.List count: %w", err)
	}
	if err := query.Preload("Owner").Offset(offset).Limit(limit).Order("created_at DESC").Find(&machines).Error; err != nil {
		return nil, 0, fmt.Errorf("machineRepository.List: %w", err)
	}
	return machines, total, nil
}

func (r *machineRepository) Update(ctx context.Context, m *models.Machine) error {
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("machineRepository.Update: %w", err)
	}
	return nil
}

func (r *machineRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Machine{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("machineRepository.Delete: %w", err)
	}
	return nil
}

func (r *machineRepository) CreateLog(ctx context.Context, l *models.MachineLog) error {
	if err := r.db.WithContext(ctx).Create(l).Error; err != nil {
		return fmt.Errorf("machineRepository.CreateLog: %w", err)
	}
	return nil
}

func (r *machineRepository) ListLogs(ctx context.Context, machineID uuid.UUID, offset, limit int) ([]models.MachineLog, int64, error) {
	var logs []models.MachineLog
	var total int64
	q := r.db.WithContext(ctx).Model(&models.MachineLog{}).Where("machine_id = ?", machineID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("machineRepository.ListLogs count: %w", err)
	}
	if err := q.Preload("Operator").Offset(offset).Limit(limit).Order("date DESC").Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("machineRepository.ListLogs: %w", err)
	}
	return logs, total, nil
}
