package services

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/google/uuid"
)

type CreateTripRequest struct {
	JobID        string  `json:"job_id"`
	TruckID      string  `json:"truck_id" validate:"required,uuid4"`
	DriverID     string  `json:"driver_id" validate:"required,uuid4"`
	Date         string  `json:"date" validate:"required"`
	FromLocation string  `json:"from_location"`
	ToLocation   string  `json:"to_location"`
	Loads        int     `json:"loads" validate:"gte=1"`
	Tonnage      float64 `json:"tonnage" validate:"gte=0"`
	RatePerTon   float64 `json:"rate_per_ton" validate:"gte=0"`
	Notes        string  `json:"notes"`
}

type TransportService interface {
	Create(ctx context.Context, req CreateTripRequest) (*models.TransportTrip, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.TransportTrip, error)
	List(ctx context.Context, jobID, truckID, driverID *uuid.UUID, from, to *time.Time, page, pageSize int) ([]models.TransportTrip, int64, error)
	Update(ctx context.Context, id uuid.UUID, req CreateTripRequest) (*models.TransportTrip, error)
	GetSummary(ctx context.Context, from, to time.Time) ([]repositories.TransportSummary, error)
}

type transportService struct {
	repo repositories.TransportRepository
}

func NewTransportService(repo repositories.TransportRepository) TransportService {
	return &transportService{repo: repo}
}

func (s *transportService) Create(ctx context.Context, req CreateTripRequest) (*models.TransportTrip, error) {
	truckID, _ := uuid.Parse(req.TruckID)
	driverID, _ := uuid.Parse(req.DriverID)
	date, _ := time.Parse("2006-01-02", req.Date)

	trip := &models.TransportTrip{
		ID:           uuid.New(),
		TruckID:      truckID,
		DriverID:     driverID,
		Date:         date,
		FromLocation: req.FromLocation,
		ToLocation:   req.ToLocation,
		Loads:        req.Loads,
		Tonnage:      req.Tonnage,
		RatePerTon:   req.RatePerTon,
		TotalAmount:  req.Tonnage * req.RatePerTon,
		Notes:        req.Notes,
	}
	if req.JobID != "" {
		jid, _ := uuid.Parse(req.JobID)
		trip.JobID = &jid
	}
	if err := s.repo.Create(ctx, trip); err != nil {
		return nil, fmt.Errorf("transportService.Create: %w", err)
	}
	return trip, nil
}

func (s *transportService) GetByID(ctx context.Context, id uuid.UUID) (*models.TransportTrip, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *transportService) List(ctx context.Context, jobID, truckID, driverID *uuid.UUID, from, to *time.Time, page, pageSize int) ([]models.TransportTrip, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, jobID, truckID, driverID, from, to, offset, pageSize)
}

func (s *transportService) Update(ctx context.Context, id uuid.UUID, req CreateTripRequest) (*models.TransportTrip, error) {
	trip, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transportService.Update find: %w", err)
	}
	if req.Loads > 0 {
		trip.Loads = req.Loads
	}
	if req.Tonnage > 0 {
		trip.Tonnage = req.Tonnage
	}
	if req.RatePerTon > 0 {
		trip.RatePerTon = req.RatePerTon
	}
	trip.TotalAmount = trip.Tonnage * trip.RatePerTon
	if req.Notes != "" {
		trip.Notes = req.Notes
	}
	if err := s.repo.Update(ctx, trip); err != nil {
		return nil, fmt.Errorf("transportService.Update: %w", err)
	}
	return trip, nil
}

func (s *transportService) GetSummary(ctx context.Context, from, to time.Time) ([]repositories.TransportSummary, error) {
	return s.repo.GetSummary(ctx, from, to)
}
