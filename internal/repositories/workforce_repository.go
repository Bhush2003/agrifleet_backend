package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkforceRepository interface {
	// Drivers
	CreateDriver(ctx context.Context, d *models.Driver) error
	FindDriverByID(ctx context.Context, id uuid.UUID) (*models.Driver, error)
	ListDrivers(ctx context.Context, offset, limit int) ([]models.Driver, int64, error)
	UpdateDriver(ctx context.Context, d *models.Driver) error

	// Workers
	CreateWorker(ctx context.Context, w *models.Worker) error
	FindWorkerByID(ctx context.Context, id uuid.UUID) (*models.Worker, error)
	ListWorkers(ctx context.Context, offset, limit int) ([]models.Worker, int64, error)
	UpdateWorker(ctx context.Context, w *models.Worker) error

	// Attendance
	CreateAttendance(ctx context.Context, a *models.Attendance) error
	BulkCreateAttendance(ctx context.Context, records []models.Attendance) error
	ListAttendance(ctx context.Context, workerID *uuid.UUID, from, to *time.Time, offset, limit int) ([]models.Attendance, int64, error)
	GetPayrollSummary(ctx context.Context, from, to time.Time) ([]PayrollSummary, error)
}

// PayrollSummary aggregates wages per worker for a period.
type PayrollSummary struct {
	WorkerID    uuid.UUID `json:"worker_id"`
	WorkerName  string    `json:"worker_name"`
	DaysPresent int       `json:"days_present"`
	HalfDays    int       `json:"half_days"`
	TotalWage   float64   `json:"total_wage"`
}

type workforceRepository struct {
	db *gorm.DB
}

func NewWorkforceRepository(db *gorm.DB) WorkforceRepository {
	return &workforceRepository{db: db}
}

func (r *workforceRepository) CreateDriver(ctx context.Context, d *models.Driver) error {
	if err := r.db.WithContext(ctx).Create(d).Error; err != nil {
		return fmt.Errorf("workforceRepository.CreateDriver: %w", err)
	}
	return nil
}

func (r *workforceRepository) FindDriverByID(ctx context.Context, id uuid.UUID) (*models.Driver, error) {
	var d models.Driver
	if err := r.db.WithContext(ctx).Preload("User").Preload("Machine").First(&d, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("workforceRepository.FindDriverByID: %w", err)
	}
	return &d, nil
}

func (r *workforceRepository) ListDrivers(ctx context.Context, offset, limit int) ([]models.Driver, int64, error) {
	var drivers []models.Driver
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Driver{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("workforceRepository.ListDrivers count: %w", err)
	}
	if err := q.Preload("User").Preload("Machine").Offset(offset).Limit(limit).Find(&drivers).Error; err != nil {
		return nil, 0, fmt.Errorf("workforceRepository.ListDrivers: %w", err)
	}
	return drivers, total, nil
}

func (r *workforceRepository) UpdateDriver(ctx context.Context, d *models.Driver) error {
	if err := r.db.WithContext(ctx).Save(d).Error; err != nil {
		return fmt.Errorf("workforceRepository.UpdateDriver: %w", err)
	}
	return nil
}

func (r *workforceRepository) CreateWorker(ctx context.Context, w *models.Worker) error {
	if err := r.db.WithContext(ctx).Create(w).Error; err != nil {
		return fmt.Errorf("workforceRepository.CreateWorker: %w", err)
	}
	return nil
}

func (r *workforceRepository) FindWorkerByID(ctx context.Context, id uuid.UUID) (*models.Worker, error) {
	var w models.Worker
	if err := r.db.WithContext(ctx).First(&w, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("workforceRepository.FindWorkerByID: %w", err)
	}
	return &w, nil
}

func (r *workforceRepository) ListWorkers(ctx context.Context, offset, limit int) ([]models.Worker, int64, error) {
	var workers []models.Worker
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Worker{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("workforceRepository.ListWorkers count: %w", err)
	}
	if err := q.Offset(offset).Limit(limit).Find(&workers).Error; err != nil {
		return nil, 0, fmt.Errorf("workforceRepository.ListWorkers: %w", err)
	}
	return workers, total, nil
}

func (r *workforceRepository) UpdateWorker(ctx context.Context, w *models.Worker) error {
	if err := r.db.WithContext(ctx).Save(w).Error; err != nil {
		return fmt.Errorf("workforceRepository.UpdateWorker: %w", err)
	}
	return nil
}

func (r *workforceRepository) CreateAttendance(ctx context.Context, a *models.Attendance) error {
	if err := r.db.WithContext(ctx).Create(a).Error; err != nil {
		return fmt.Errorf("workforceRepository.CreateAttendance: %w", err)
	}
	return nil
}

func (r *workforceRepository) BulkCreateAttendance(ctx context.Context, records []models.Attendance) error {
	if err := r.db.WithContext(ctx).CreateInBatches(records, 100).Error; err != nil {
		return fmt.Errorf("workforceRepository.BulkCreateAttendance: %w", err)
	}
	return nil
}

func (r *workforceRepository) ListAttendance(ctx context.Context, workerID *uuid.UUID, from, to *time.Time, offset, limit int) ([]models.Attendance, int64, error) {
	var records []models.Attendance
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Attendance{})
	if workerID != nil {
		q = q.Where("worker_id = ?", *workerID)
	}
	if from != nil {
		q = q.Where("date >= ?", *from)
	}
	if to != nil {
		q = q.Where("date <= ?", *to)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("workforceRepository.ListAttendance count: %w", err)
	}
	if err := q.Preload("Worker").Offset(offset).Limit(limit).Order("date DESC").Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("workforceRepository.ListAttendance: %w", err)
	}
	return records, total, nil
}

func (r *workforceRepository) GetPayrollSummary(ctx context.Context, from, to time.Time) ([]PayrollSummary, error) {
	var results []PayrollSummary
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			a.worker_id,
			w.name AS worker_name,
			COUNT(CASE WHEN a.status = 'present' THEN 1 END) AS days_present,
			COUNT(CASE WHEN a.status = 'half_day' THEN 1 END) AS half_days,
			SUM(a.wage_earned) AS total_wage
		FROM attendances a
		JOIN workers w ON w.id = a.worker_id
		WHERE a.date BETWEEN ? AND ?
		  AND w.deleted_at IS NULL
		GROUP BY a.worker_id, w.name
		ORDER BY w.name
	`, from, to).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("workforceRepository.GetPayrollSummary: %w", err)
	}
	return results, nil
}
