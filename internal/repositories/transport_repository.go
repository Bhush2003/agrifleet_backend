package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransportRepository interface {
	Create(ctx context.Context, t *models.TransportTrip) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.TransportTrip, error)
	List(ctx context.Context, jobID, truckID, driverID *uuid.UUID, from, to *time.Time, offset, limit int) ([]models.TransportTrip, int64, error)
	Update(ctx context.Context, t *models.TransportTrip) error
	GetSummary(ctx context.Context, from, to time.Time) ([]TransportSummary, error)
}

// TransportSummary aggregates trips per driver.
type TransportSummary struct {
	DriverID    uuid.UUID `json:"driver_id"`
	DriverName  string    `json:"driver_name"`
	TotalTrips  int       `json:"total_trips"`
	TotalLoads  int       `json:"total_loads"`
	TotalTonnage float64  `json:"total_tonnage"`
	TotalAmount float64   `json:"total_amount"`
}

type transportRepository struct {
	db *gorm.DB
}

func NewTransportRepository(db *gorm.DB) TransportRepository {
	return &transportRepository{db: db}
}

func (r *transportRepository) Create(ctx context.Context, t *models.TransportTrip) error {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return fmt.Errorf("transportRepository.Create: %w", err)
	}
	return nil
}

func (r *transportRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.TransportTrip, error) {
	var t models.TransportTrip
	if err := r.db.WithContext(ctx).Preload("Truck").Preload("Driver").Preload("Job").First(&t, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("transportRepository.FindByID: %w", err)
	}
	return &t, nil
}

func (r *transportRepository) List(ctx context.Context, jobID, truckID, driverID *uuid.UUID, from, to *time.Time, offset, limit int) ([]models.TransportTrip, int64, error) {
	var trips []models.TransportTrip
	var total int64
	q := r.db.WithContext(ctx).Model(&models.TransportTrip{})
	if jobID != nil {
		q = q.Where("job_id = ?", *jobID)
	}
	if truckID != nil {
		q = q.Where("truck_id = ?", *truckID)
	}
	if driverID != nil {
		q = q.Where("driver_id = ?", *driverID)
	}
	if from != nil {
		q = q.Where("date >= ?", *from)
	}
	if to != nil {
		q = q.Where("date <= ?", *to)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("transportRepository.List count: %w", err)
	}
	if err := q.Preload("Truck").Preload("Driver.User").Offset(offset).Limit(limit).Order("date DESC").Find(&trips).Error; err != nil {
		return nil, 0, fmt.Errorf("transportRepository.List: %w", err)
	}
	return trips, total, nil
}

func (r *transportRepository) Update(ctx context.Context, t *models.TransportTrip) error {
	if err := r.db.WithContext(ctx).Save(t).Error; err != nil {
		return fmt.Errorf("transportRepository.Update: %w", err)
	}
	return nil
}

func (r *transportRepository) GetSummary(ctx context.Context, from, to time.Time) ([]TransportSummary, error) {
	var results []TransportSummary
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			t.driver_id,
			u.name AS driver_name,
			COUNT(*) AS total_trips,
			SUM(t.loads) AS total_loads,
			SUM(t.tonnage) AS total_tonnage,
			SUM(t.total_amount) AS total_amount
		FROM transport_trips t
		JOIN drivers d ON d.id = t.driver_id
		JOIN users u ON u.id = d.user_id
		WHERE t.date BETWEEN ? AND ?
		  AND t.deleted_at IS NULL
		GROUP BY t.driver_id, u.name
		ORDER BY total_amount DESC
	`, from, to).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("transportRepository.GetSummary: %w", err)
	}
	return results, nil
}
