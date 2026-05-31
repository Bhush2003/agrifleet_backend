package services

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/google/uuid"
)

type CreateDriverRequest struct {
	UserID        string `json:"user_id" validate:"required,uuid4"`
	LicenseNumber string `json:"license_number" validate:"required"`
	LicenseExpiry string `json:"license_expiry" validate:"required"`
	MachineID     string `json:"machine_id"`
}

type CreateWorkerRequest struct {
	Name      string  `json:"name" validate:"required,min=2"`
	Phone     string  `json:"phone"`
	Role      string  `json:"role"`
	DailyWage float64 `json:"daily_wage" validate:"gte=0"`
	JoinDate  string  `json:"join_date"`
}

type AttendanceRecord struct {
	WorkerID string `json:"worker_id" validate:"required,uuid4"`
	Date     string `json:"date" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=present absent half_day on_leave"`
	Notes    string `json:"notes"`
}

type BulkAttendanceRequest struct {
	Records []AttendanceRecord `json:"records" validate:"required,min=1"`
}

type WorkforceService interface {
	CreateDriver(ctx context.Context, req CreateDriverRequest) (*models.Driver, error)
	GetDriver(ctx context.Context, id uuid.UUID) (*models.Driver, error)
	ListDrivers(ctx context.Context, page, pageSize int) ([]models.Driver, int64, error)
	UpdateDriver(ctx context.Context, id uuid.UUID, req CreateDriverRequest) (*models.Driver, error)

	CreateWorker(ctx context.Context, req CreateWorkerRequest) (*models.Worker, error)
	GetWorker(ctx context.Context, id uuid.UUID) (*models.Worker, error)
	ListWorkers(ctx context.Context, page, pageSize int) ([]models.Worker, int64, error)
	UpdateWorker(ctx context.Context, id uuid.UUID, req CreateWorkerRequest) (*models.Worker, error)

	MarkAttendance(ctx context.Context, req BulkAttendanceRequest) error
	ListAttendance(ctx context.Context, workerID *uuid.UUID, from, to *time.Time, page, pageSize int) ([]models.Attendance, int64, error)
	GetPayroll(ctx context.Context, from, to time.Time) ([]repositories.PayrollSummary, error)
}

type workforceService struct {
	repo repositories.WorkforceRepository
}

func NewWorkforceService(repo repositories.WorkforceRepository) WorkforceService {
	return &workforceService{repo: repo}
}

func (s *workforceService) CreateDriver(ctx context.Context, req CreateDriverRequest) (*models.Driver, error) {
	userID, _ := uuid.Parse(req.UserID)
	expiry, _ := time.Parse("2006-01-02", req.LicenseExpiry)

	d := &models.Driver{
		ID:            uuid.New(),
		UserID:        userID,
		LicenseNumber: req.LicenseNumber,
		LicenseExpiry: expiry,
		Status:        models.DriverStatusActive,
	}
	if req.MachineID != "" {
		mid, _ := uuid.Parse(req.MachineID)
		d.MachineID = &mid
	}
	if err := s.repo.CreateDriver(ctx, d); err != nil {
		return nil, fmt.Errorf("workforceService.CreateDriver: %w", err)
	}
	return d, nil
}

func (s *workforceService) GetDriver(ctx context.Context, id uuid.UUID) (*models.Driver, error) {
	d, err := s.repo.FindDriverByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("workforceService.GetDriver: %w", err)
	}
	return d, nil
}

func (s *workforceService) ListDrivers(ctx context.Context, page, pageSize int) ([]models.Driver, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListDrivers(ctx, offset, pageSize)
}

func (s *workforceService) UpdateDriver(ctx context.Context, id uuid.UUID, req CreateDriverRequest) (*models.Driver, error) {
	d, err := s.repo.FindDriverByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("workforceService.UpdateDriver find: %w", err)
	}
	if req.LicenseNumber != "" {
		d.LicenseNumber = req.LicenseNumber
	}
	if req.LicenseExpiry != "" {
		expiry, _ := time.Parse("2006-01-02", req.LicenseExpiry)
		d.LicenseExpiry = expiry
	}
	if req.MachineID != "" {
		mid, _ := uuid.Parse(req.MachineID)
		d.MachineID = &mid
	}
	if err := s.repo.UpdateDriver(ctx, d); err != nil {
		return nil, fmt.Errorf("workforceService.UpdateDriver: %w", err)
	}
	return d, nil
}

func (s *workforceService) CreateWorker(ctx context.Context, req CreateWorkerRequest) (*models.Worker, error) {
	joinDate := time.Now()
	if req.JoinDate != "" {
		joinDate, _ = time.Parse("2006-01-02", req.JoinDate)
	}
	w := &models.Worker{
		ID:        uuid.New(),
		Name:      req.Name,
		Phone:     req.Phone,
		Role:      req.Role,
		DailyWage: req.DailyWage,
		JoinDate:  joinDate,
		Status:    models.WorkerStatusActive,
	}
	if err := s.repo.CreateWorker(ctx, w); err != nil {
		return nil, fmt.Errorf("workforceService.CreateWorker: %w", err)
	}
	return w, nil
}

func (s *workforceService) GetWorker(ctx context.Context, id uuid.UUID) (*models.Worker, error) {
	return s.repo.FindWorkerByID(ctx, id)
}

func (s *workforceService) ListWorkers(ctx context.Context, page, pageSize int) ([]models.Worker, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListWorkers(ctx, offset, pageSize)
}

func (s *workforceService) UpdateWorker(ctx context.Context, id uuid.UUID, req CreateWorkerRequest) (*models.Worker, error) {
	w, err := s.repo.FindWorkerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("workforceService.UpdateWorker find: %w", err)
	}
	if req.Name != "" {
		w.Name = req.Name
	}
	if req.Phone != "" {
		w.Phone = req.Phone
	}
	if req.Role != "" {
		w.Role = req.Role
	}
	if req.DailyWage > 0 {
		w.DailyWage = req.DailyWage
	}
	if err := s.repo.UpdateWorker(ctx, w); err != nil {
		return nil, fmt.Errorf("workforceService.UpdateWorker: %w", err)
	}
	return w, nil
}

func (s *workforceService) MarkAttendance(ctx context.Context, req BulkAttendanceRequest) error {
	var records []models.Attendance
	for _, r := range req.Records {
		workerID, _ := uuid.Parse(r.WorkerID)
		date, _ := time.Parse("2006-01-02", r.Date)

		worker, err := s.repo.FindWorkerByID(ctx, workerID)
		if err != nil {
			return fmt.Errorf("workforceService.MarkAttendance worker %s: %w", r.WorkerID, err)
		}

		wage := 0.0
		switch models.AttendanceStatus(r.Status) {
		case models.AttendancePresent:
			wage = worker.DailyWage
		case models.AttendanceHalfDay:
			wage = worker.DailyWage / 2
		}

		records = append(records, models.Attendance{
			ID:         uuid.New(),
			WorkerID:   workerID,
			Date:       date,
			Status:     models.AttendanceStatus(r.Status),
			WageEarned: wage,
			Notes:      r.Notes,
		})
	}
	if err := s.repo.BulkCreateAttendance(ctx, records); err != nil {
		return fmt.Errorf("workforceService.MarkAttendance: %w", err)
	}
	return nil
}

func (s *workforceService) ListAttendance(ctx context.Context, workerID *uuid.UUID, from, to *time.Time, page, pageSize int) ([]models.Attendance, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListAttendance(ctx, workerID, from, to, offset, pageSize)
}

func (s *workforceService) GetPayroll(ctx context.Context, from, to time.Time) ([]repositories.PayrollSummary, error) {
	return s.repo.GetPayrollSummary(ctx, from, to)
}
