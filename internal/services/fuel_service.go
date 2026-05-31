package services

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/google/uuid"
)

type CreateFuelEntryRequest struct {
	MachineID    string  `json:"machine_id" validate:"required,uuid4"`
	Date         string  `json:"date" validate:"required"`
	Liters       float64 `json:"liters" validate:"required,gt=0"`
	CostPerLiter float64 `json:"cost_per_liter" validate:"required,gt=0"`
	Odometer     float64 `json:"odometer"`
	Notes        string  `json:"notes"`
}

type FuelService interface {
	Create(ctx context.Context, filledBy uuid.UUID, req CreateFuelEntryRequest) (*models.FuelEntry, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.FuelEntry, error)
	List(ctx context.Context, machineID *uuid.UUID, from, to *time.Time, page, pageSize int) ([]models.FuelEntry, int64, error)
	Update(ctx context.Context, id uuid.UUID, req CreateFuelEntryRequest) (*models.FuelEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAnalytics(ctx context.Context, from, to time.Time) ([]repositories.FuelAnalytics, error)
}

type fuelService struct {
	repo repositories.FuelRepository
}

func NewFuelService(repo repositories.FuelRepository) FuelService {
	return &fuelService{repo: repo}
}

func (s *fuelService) Create(ctx context.Context, filledBy uuid.UUID, req CreateFuelEntryRequest) (*models.FuelEntry, error) {
	machineID, _ := uuid.Parse(req.MachineID)
	date, _ := time.Parse("2006-01-02", req.Date)

	entry := &models.FuelEntry{
		ID:           uuid.New(),
		MachineID:    machineID,
		Date:         date,
		Liters:       req.Liters,
		CostPerLiter: req.CostPerLiter,
		TotalCost:    req.Liters * req.CostPerLiter,
		FilledBy:     filledBy,
		Odometer:     req.Odometer,
		Notes:        req.Notes,
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("fuelService.Create: %w", err)
	}
	return entry, nil
}

func (s *fuelService) GetByID(ctx context.Context, id uuid.UUID) (*models.FuelEntry, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *fuelService) List(ctx context.Context, machineID *uuid.UUID, from, to *time.Time, page, pageSize int) ([]models.FuelEntry, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, machineID, from, to, offset, pageSize)
}

func (s *fuelService) Update(ctx context.Context, id uuid.UUID, req CreateFuelEntryRequest) (*models.FuelEntry, error) {
	entry, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("fuelService.Update find: %w", err)
	}
	if req.Liters > 0 {
		entry.Liters = req.Liters
	}
	if req.CostPerLiter > 0 {
		entry.CostPerLiter = req.CostPerLiter
	}
	entry.TotalCost = entry.Liters * entry.CostPerLiter
	if req.Odometer > 0 {
		entry.Odometer = req.Odometer
	}
	if req.Notes != "" {
		entry.Notes = req.Notes
	}
	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, fmt.Errorf("fuelService.Update: %w", err)
	}
	return entry, nil
}

func (s *fuelService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *fuelService) GetAnalytics(ctx context.Context, from, to time.Time) ([]repositories.FuelAnalytics, error) {
	return s.repo.GetAnalytics(ctx, from, to)
}
