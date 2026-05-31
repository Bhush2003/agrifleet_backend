package repositories

import (
	"context"
	"fmt"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HarvestingRepository interface {
	CreateJob(ctx context.Context, j *models.HarvestingJob) error
	FindJobByID(ctx context.Context, id uuid.UUID) (*models.HarvestingJob, error)
	ListJobs(ctx context.Context, status *models.JobStatus, offset, limit int) ([]models.HarvestingJob, int64, error)
	UpdateJob(ctx context.Context, j *models.HarvestingJob) error
	CreateLog(ctx context.Context, l *models.HarvestingLog) error
	ListLogs(ctx context.Context, jobID uuid.UUID, offset, limit int) ([]models.HarvestingLog, int64, error)
	GetJobSummary(ctx context.Context, jobID uuid.UUID) (*HarvestingJobSummary, error)
}

// HarvestingJobSummary aggregates totals for a job.
type HarvestingJobSummary struct {
	JobID       uuid.UUID `json:"job_id"`
	TotalAcres  float64   `json:"total_acres"`
	TotalTonnage float64  `json:"total_tonnage"`
	DaysWorked  int       `json:"days_worked"`
}

type harvestingRepository struct {
	db *gorm.DB
}

func NewHarvestingRepository(db *gorm.DB) HarvestingRepository {
	return &harvestingRepository{db: db}
}

func (r *harvestingRepository) CreateJob(ctx context.Context, j *models.HarvestingJob) error {
	if err := r.db.WithContext(ctx).Create(j).Error; err != nil {
		return fmt.Errorf("harvestingRepository.CreateJob: %w", err)
	}
	return nil
}

func (r *harvestingRepository) FindJobByID(ctx context.Context, id uuid.UUID) (*models.HarvestingJob, error) {
	var j models.HarvestingJob
	if err := r.db.WithContext(ctx).Preload("Machine").Preload("Project").First(&j, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("harvestingRepository.FindJobByID: %w", err)
	}
	return &j, nil
}

func (r *harvestingRepository) ListJobs(ctx context.Context, status *models.JobStatus, offset, limit int) ([]models.HarvestingJob, int64, error) {
	var jobs []models.HarvestingJob
	var total int64
	q := r.db.WithContext(ctx).Model(&models.HarvestingJob{})
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("harvestingRepository.ListJobs count: %w", err)
	}
	if err := q.Preload("Machine").Offset(offset).Limit(limit).Order("start_date DESC").Find(&jobs).Error; err != nil {
		return nil, 0, fmt.Errorf("harvestingRepository.ListJobs: %w", err)
	}
	return jobs, total, nil
}

func (r *harvestingRepository) UpdateJob(ctx context.Context, j *models.HarvestingJob) error {
	if err := r.db.WithContext(ctx).Save(j).Error; err != nil {
		return fmt.Errorf("harvestingRepository.UpdateJob: %w", err)
	}
	return nil
}

func (r *harvestingRepository) CreateLog(ctx context.Context, l *models.HarvestingLog) error {
	if err := r.db.WithContext(ctx).Create(l).Error; err != nil {
		return fmt.Errorf("harvestingRepository.CreateLog: %w", err)
	}
	return nil
}

func (r *harvestingRepository) ListLogs(ctx context.Context, jobID uuid.UUID, offset, limit int) ([]models.HarvestingLog, int64, error) {
	var logs []models.HarvestingLog
	var total int64
	q := r.db.WithContext(ctx).Model(&models.HarvestingLog{}).Where("job_id = ?", jobID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("harvestingRepository.ListLogs count: %w", err)
	}
	if err := q.Preload("Machine").Preload("Driver").Offset(offset).Limit(limit).Order("date DESC").Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("harvestingRepository.ListLogs: %w", err)
	}
	return logs, total, nil
}

func (r *harvestingRepository) GetJobSummary(ctx context.Context, jobID uuid.UUID) (*HarvestingJobSummary, error) {
	var summary HarvestingJobSummary
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			job_id,
			SUM(acres_done) AS total_acres,
			SUM(tonnage) AS total_tonnage,
			COUNT(DISTINCT date) AS days_worked
		FROM harvesting_logs
		WHERE job_id = ?
		GROUP BY job_id
	`, jobID).Scan(&summary).Error
	if err != nil {
		return nil, fmt.Errorf("harvestingRepository.GetJobSummary: %w", err)
	}
	summary.JobID = jobID
	return &summary, nil
}
