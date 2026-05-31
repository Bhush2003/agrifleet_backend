package services

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/google/uuid"
)

type CreateMaintenanceRequest struct {
	MachineID   string `json:"machine_id" validate:"required,uuid4"`
	Type        string `json:"type" validate:"required,oneof=scheduled breakdown preventive"`
	Date        string `json:"date" validate:"required"`
	Description string `json:"description" validate:"required"`
	Cost        float64 `json:"cost"`
	NextDueDate string `json:"next_due_date"`
	Technician  string `json:"technician"`
	Notes       string `json:"notes"`
}

type MaintenanceService interface {
	Create(ctx context.Context, req CreateMaintenanceRequest) (*models.Maintenance, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Maintenance, error)
	List(ctx context.Context, machineID *uuid.UUID, page, pageSize int) ([]models.Maintenance, int64, error)
	Update(ctx context.Context, id uuid.UUID, req CreateMaintenanceRequest) (*models.Maintenance, error)
	ListUpcoming(ctx context.Context, days int) ([]models.Maintenance, error)
	ListOverdue(ctx context.Context) ([]models.Maintenance, error)
}

type maintenanceService struct {
	repo repositories.MaintenanceRepository
}

func NewMaintenanceService(repo repositories.MaintenanceRepository) MaintenanceService {
	return &maintenanceService{repo: repo}
}

func (s *maintenanceService) Create(ctx context.Context, req CreateMaintenanceRequest) (*models.Maintenance, error) {
	machineID, _ := uuid.Parse(req.MachineID)
	date, _ := time.Parse("2006-01-02", req.Date)

	m := &models.Maintenance{
		ID:          uuid.New(),
		MachineID:   machineID,
		Type:        models.MaintenanceType(req.Type),
		Date:        date,
		Description: req.Description,
		Cost:        req.Cost,
		Status:      models.MaintenancePending,
		Technician:  req.Technician,
		Notes:       req.Notes,
	}
	if req.NextDueDate != "" {
		t, _ := time.Parse("2006-01-02", req.NextDueDate)
		m.NextDueDate = &t
	}
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("maintenanceService.Create: %w", err)
	}
	return m, nil
}

func (s *maintenanceService) GetByID(ctx context.Context, id uuid.UUID) (*models.Maintenance, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *maintenanceService) List(ctx context.Context, machineID *uuid.UUID, page, pageSize int) ([]models.Maintenance, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, machineID, offset, pageSize)
}

func (s *maintenanceService) Update(ctx context.Context, id uuid.UUID, req CreateMaintenanceRequest) (*models.Maintenance, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("maintenanceService.Update find: %w", err)
	}
	if req.Description != "" {
		m.Description = req.Description
	}
	if req.Cost > 0 {
		m.Cost = req.Cost
	}
	if req.Technician != "" {
		m.Technician = req.Technician
	}
	if req.Notes != "" {
		m.Notes = req.Notes
	}
	if req.NextDueDate != "" {
		t, _ := time.Parse("2006-01-02", req.NextDueDate)
		m.NextDueDate = &t
	}
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, fmt.Errorf("maintenanceService.Update: %w", err)
	}
	return m, nil
}

func (s *maintenanceService) ListUpcoming(ctx context.Context, days int) ([]models.Maintenance, error) {
	return s.repo.ListUpcoming(ctx, days)
}

func (s *maintenanceService) ListOverdue(ctx context.Context) ([]models.Maintenance, error) {
	return s.repo.ListOverdue(ctx)
}
