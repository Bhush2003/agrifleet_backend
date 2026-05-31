package services

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/google/uuid"
)

type CreateHarvestingJobRequest struct {
	FarmName   string `json:"farm_name" validate:"required"`
	Location   string `json:"location"`
	CropType   string `json:"crop_type"`
	StartDate  string `json:"start_date" validate:"required"`
	TotalAcres float64 `json:"total_acres"`
	MachineID  string `json:"machine_id"`
	ProjectID  string `json:"project_id"`
	Notes      string `json:"notes"`
}

type AddHarvestingLogRequest struct {
	Date      string  `json:"date" validate:"required"`
	MachineID string  `json:"machine_id" validate:"required,uuid4"`
	DriverID  string  `json:"driver_id"`
	AcresDone float64 `json:"acres_done" validate:"gte=0"`
	Tonnage   float64 `json:"tonnage" validate:"gte=0"`
	Notes     string  `json:"notes"`
}

type HarvestingService interface {
	CreateJob(ctx context.Context, req CreateHarvestingJobRequest) (*models.HarvestingJob, error)
	GetJob(ctx context.Context, id uuid.UUID) (*models.HarvestingJob, error)
	ListJobs(ctx context.Context, status *models.JobStatus, page, pageSize int) ([]models.HarvestingJob, int64, error)
	UpdateJob(ctx context.Context, id uuid.UUID, req CreateHarvestingJobRequest) (*models.HarvestingJob, error)
	AddLog(ctx context.Context, jobID uuid.UUID, req AddHarvestingLogRequest) (*models.HarvestingLog, error)
	ListLogs(ctx context.Context, jobID uuid.UUID, page, pageSize int) ([]models.HarvestingLog, int64, error)
	GetJobSummary(ctx context.Context, jobID uuid.UUID) (*repositories.HarvestingJobSummary, error)
}

type harvestingService struct {
	repo repositories.HarvestingRepository
}

func NewHarvestingService(repo repositories.HarvestingRepository) HarvestingService {
	return &harvestingService{repo: repo}
}

func (s *harvestingService) CreateJob(ctx context.Context, req CreateHarvestingJobRequest) (*models.HarvestingJob, error) {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	job := &models.HarvestingJob{
		ID:         uuid.New(),
		FarmName:   req.FarmName,
		Location:   req.Location,
		CropType:   req.CropType,
		StartDate:  startDate,
		TotalAcres: req.TotalAcres,
		Status:     models.JobUpcoming,
		Notes:      req.Notes,
	}
	if req.MachineID != "" {
		mid, _ := uuid.Parse(req.MachineID)
		job.MachineID = &mid
	}
	if req.ProjectID != "" {
		pid, _ := uuid.Parse(req.ProjectID)
		job.ProjectID = &pid
	}
	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("harvestingService.CreateJob: %w", err)
	}
	return job, nil
}

func (s *harvestingService) GetJob(ctx context.Context, id uuid.UUID) (*models.HarvestingJob, error) {
	return s.repo.FindJobByID(ctx, id)
}

func (s *harvestingService) ListJobs(ctx context.Context, status *models.JobStatus, page, pageSize int) ([]models.HarvestingJob, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListJobs(ctx, status, offset, pageSize)
}

func (s *harvestingService) UpdateJob(ctx context.Context, id uuid.UUID, req CreateHarvestingJobRequest) (*models.HarvestingJob, error) {
	job, err := s.repo.FindJobByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("harvestingService.UpdateJob find: %w", err)
	}
	if req.FarmName != "" {
		job.FarmName = req.FarmName
	}
	if req.Location != "" {
		job.Location = req.Location
	}
	if req.CropType != "" {
		job.CropType = req.CropType
	}
	if req.TotalAcres > 0 {
		job.TotalAcres = req.TotalAcres
	}
	if req.Notes != "" {
		job.Notes = req.Notes
	}
	if err := s.repo.UpdateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("harvestingService.UpdateJob: %w", err)
	}
	return job, nil
}

func (s *harvestingService) AddLog(ctx context.Context, jobID uuid.UUID, req AddHarvestingLogRequest) (*models.HarvestingLog, error) {
	machineID, _ := uuid.Parse(req.MachineID)
	date, _ := time.Parse("2006-01-02", req.Date)

	log := &models.HarvestingLog{
		ID:        uuid.New(),
		JobID:     jobID,
		Date:      date,
		MachineID: machineID,
		AcresDone: req.AcresDone,
		Tonnage:   req.Tonnage,
		Notes:     req.Notes,
	}
	if req.DriverID != "" {
		did, _ := uuid.Parse(req.DriverID)
		log.DriverID = did
	}
	if err := s.repo.CreateLog(ctx, log); err != nil {
		return nil, fmt.Errorf("harvestingService.AddLog: %w", err)
	}
	return log, nil
}

func (s *harvestingService) ListLogs(ctx context.Context, jobID uuid.UUID, page, pageSize int) ([]models.HarvestingLog, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListLogs(ctx, jobID, offset, pageSize)
}

func (s *harvestingService) GetJobSummary(ctx context.Context, jobID uuid.UUID) (*repositories.HarvestingJobSummary, error) {
	return s.repo.GetJobSummary(ctx, jobID)
}
