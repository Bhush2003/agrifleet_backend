package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FuelRepository interface {
	Create(ctx context.Context, f *models.FuelEntry) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.FuelEntry, error)
	List(ctx context.Context, machineID *uuid.UUID, from, to *time.Time, offset, limit int) ([]models.FuelEntry, int64, error)
	Update(ctx context.Context, f *models.FuelEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAnalytics(ctx context.Context, from, to time.Time) ([]FuelAnalytics, error)
}

// FuelAnalytics aggregates fuel consumption per machine.
type FuelAnalytics struct {
	MachineID   uuid.UUID `json:"machine_id"`
	MachineName string    `json:"machine_name"`
	TotalLiters float64   `json:"total_liters"`
	TotalCost   float64   `json:"total_cost"`
	FillCount   int       `json:"fill_count"`
}

type fuelRepository struct {
	db *gorm.DB
}

func NewFuelRepository(db *gorm.DB) FuelRepository {
	return &fuelRepository{db: db}
}

func (r *fuelRepository) Create(ctx context.Context, f *models.FuelEntry) error {
	if err := r.db.WithContext(ctx).Create(f).Error; err != nil {
		return fmt.Errorf("fuelRepository.Create: %w", err)
	}
	return nil
}

func (r *fuelRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.FuelEntry, error) {
	var f models.FuelEntry
	if err := r.db.WithContext(ctx).Preload("Machine").First(&f, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("fuelRepository.FindByID: %w", err)
	}
	return &f, nil
}

func (r *fuelRepository) List(ctx context.Context, machineID *uuid.UUID, from, to *time.Time, offset, limit int) ([]models.FuelEntry, int64, error) {
	var entries []models.FuelEntry
	var total int64
	q := r.db.WithContext(ctx).Model(&models.FuelEntry{})
	if machineID != nil {
		q = q.Where("machine_id = ?", *machineID)
	}
	if from != nil {
		q = q.Where("date >= ?", *from)
	}
	if to != nil {
		q = q.Where("date <= ?", *to)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("fuelRepository.List count: %w", err)
	}
	if err := q.Preload("Machine").Offset(offset).Limit(limit).Order("date DESC").Find(&entries).Error; err != nil {
		return nil, 0, fmt.Errorf("fuelRepository.List: %w", err)
	}
	return entries, total, nil
}

func (r *fuelRepository) Update(ctx context.Context, f *models.FuelEntry) error {
	if err := r.db.WithContext(ctx).Save(f).Error; err != nil {
		return fmt.Errorf("fuelRepository.Update: %w", err)
	}
	return nil
}

func (r *fuelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.FuelEntry{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("fuelRepository.Delete: %w", err)
	}
	return nil
}

func (r *fuelRepository) GetAnalytics(ctx context.Context, from, to time.Time) ([]FuelAnalytics, error) {
	var results []FuelAnalytics
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			f.machine_id,
			m.name AS machine_name,
			SUM(f.liters) AS total_liters,
			SUM(f.total_cost) AS total_cost,
			COUNT(*) AS fill_count
		FROM fuel_entries f
		JOIN machines m ON m.id = f.machine_id
		WHERE f.date BETWEEN ? AND ?
		  AND f.deleted_at IS NULL
		GROUP BY f.machine_id, m.name
		ORDER BY total_cost DESC
	`, from, to).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("fuelRepository.GetAnalytics: %w", err)
	}
	return results, nil
}
