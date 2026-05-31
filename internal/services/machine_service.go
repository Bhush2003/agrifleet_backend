package services

import (
	"context"
	"fmt"

	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/google/uuid"
)

type CreateMachineRequest struct {
	Name            string             `json:"name" validate:"required,min=2,max=100"`
	Type            models.MachineType `json:"type" validate:"required,oneof=harvester tractor truck loader other"`
	Model           string             `json:"model"`
	Year            int                `json:"year"`
	RegNumber       string             `json:"reg_number" validate:"required"`
	CurrentLocation string             `json:"current_location"`
	Notes           string             `json:"notes"`
}

type UpdateMachineRequest struct {
	Name            string                `json:"name"`
	Status          models.MachineStatus  `json:"status"`
	CurrentLocation string                `json:"current_location"`
	Notes           string                `json:"notes"`
}

type AddMachineLogRequest struct {
	Date         string  `json:"date" validate:"required"`
	WorkingHours float64 `json:"working_hours" validate:"gte=0"`
	Location     string  `json:"location"`
	OperatorID   string  `json:"operator_id"`
	Notes        string  `json:"notes"`
}

type MachineService interface {
	Create(ctx context.Context, ownerID uuid.UUID, req CreateMachineRequest) (*models.Machine, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Machine, error)
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]models.Machine, int64, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateMachineRequest) (*models.Machine, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AddLog(ctx context.Context, machineID uuid.UUID, operatorID uuid.UUID, req AddMachineLogRequest) (*models.MachineLog, error)
	ListLogs(ctx context.Context, machineID uuid.UUID, page, pageSize int) ([]models.MachineLog, int64, error)
}

type machineService struct {
	repo repositories.MachineRepository
}

func NewMachineService(repo repositories.MachineRepository) MachineService {
	return &machineService{repo: repo}
}

func (s *machineService) Create(ctx context.Context, ownerID uuid.UUID, req CreateMachineRequest) (*models.Machine, error) {
	m := &models.Machine{
		ID:              uuid.New(),
		Name:            req.Name,
		Type:            req.Type,
		Model:           req.Model,
		Year:            req.Year,
		RegNumber:       req.RegNumber,
		Status:          models.MachineStatusIdle,
		CurrentLocation: req.CurrentLocation,
		OwnerID:         ownerID,
		Notes:           req.Notes,
	}
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, fmt.Errorf("machineService.Create: %w", err)
	}
	return m, nil
}

func (s *machineService) GetByID(ctx context.Context, id uuid.UUID) (*models.Machine, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("machineService.GetByID: %w", err)
	}
	return m, nil
}

func (s *machineService) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]models.Machine, int64, error) {
	offset := (page - 1) * pageSize
	machines, total, err := s.repo.List(ctx, filters, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("machineService.List: %w", err)
	}
	return machines, total, nil
}

func (s *machineService) Update(ctx context.Context, id uuid.UUID, req UpdateMachineRequest) (*models.Machine, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("machineService.Update find: %w", err)
	}
	if req.Name != "" {
		m.Name = req.Name
	}
	if req.Status != "" {
		m.Status = req.Status
	}
	if req.CurrentLocation != "" {
		m.CurrentLocation = req.CurrentLocation
	}
	if req.Notes != "" {
		m.Notes = req.Notes
	}
	if err := s.repo.Update(ctx, m); err != nil {
		return nil, fmt.Errorf("machineService.Update: %w", err)
	}
	return m, nil
}

func (s *machineService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("machineService.Delete: %w", err)
	}
	return nil
}

func (s *machineService) AddLog(ctx context.Context, machineID uuid.UUID, operatorID uuid.UUID, req AddMachineLogRequest) (*models.MachineLog, error) {
	log := &models.MachineLog{
		ID:           uuid.New(),
		MachineID:    machineID,
		WorkingHours: req.WorkingHours,
		Location:     req.Location,
		OperatorID:   operatorID,
		Notes:        req.Notes,
	}
	if err := s.repo.CreateLog(ctx, log); err != nil {
		return nil, fmt.Errorf("machineService.AddLog: %w", err)
	}
	return log, nil
}

func (s *machineService) ListLogs(ctx context.Context, machineID uuid.UUID, page, pageSize int) ([]models.MachineLog, int64, error) {
	offset := (page - 1) * pageSize
	logs, total, err := s.repo.ListLogs(ctx, machineID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("machineService.ListLogs: %w", err)
	}
	return logs, total, nil
}
