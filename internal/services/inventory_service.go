package services

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/google/uuid"
)

type CreateSparePartRequest struct {
	Name          string  `json:"name" validate:"required,min=2"`
	Category      string  `json:"category"`
	Unit          string  `json:"unit"`
	MinStockLevel float64 `json:"min_stock_level"`
	UnitCost      float64 `json:"unit_cost"`
}

type StockAdjustRequest struct {
	Quantity      float64 `json:"quantity" validate:"required,gt=0"`
	ReferenceID   string  `json:"reference_id"`
	ReferenceType string  `json:"reference_type"`
	Notes         string  `json:"notes"`
}

type InventoryService interface {
	CreatePart(ctx context.Context, req CreateSparePartRequest) (*models.SparePart, error)
	GetPart(ctx context.Context, id uuid.UUID) (*models.SparePart, error)
	ListParts(ctx context.Context, page, pageSize int) ([]models.SparePart, int64, error)
	UpdatePart(ctx context.Context, id uuid.UUID, req CreateSparePartRequest) (*models.SparePart, error)
	StockIn(ctx context.Context, partID uuid.UUID, userID uuid.UUID, req StockAdjustRequest) error
	StockOut(ctx context.Context, partID uuid.UUID, userID uuid.UUID, req StockAdjustRequest) error
	ListLowStock(ctx context.Context) ([]models.SparePart, error)
	ListMovements(ctx context.Context, partID *uuid.UUID, page, pageSize int) ([]models.StockMovement, int64, error)
}

type inventoryService struct {
	repo repositories.InventoryRepository
}

func NewInventoryService(repo repositories.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) CreatePart(ctx context.Context, req CreateSparePartRequest) (*models.SparePart, error) {
	p := &models.SparePart{
		ID:            uuid.New(),
		Name:          req.Name,
		Category:      req.Category,
		Unit:          req.Unit,
		MinStockLevel: req.MinStockLevel,
		UnitCost:      req.UnitCost,
	}
	if err := s.repo.CreatePart(ctx, p); err != nil {
		return nil, fmt.Errorf("inventoryService.CreatePart: %w", err)
	}
	return p, nil
}

func (s *inventoryService) GetPart(ctx context.Context, id uuid.UUID) (*models.SparePart, error) {
	return s.repo.FindPartByID(ctx, id)
}

func (s *inventoryService) ListParts(ctx context.Context, page, pageSize int) ([]models.SparePart, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListParts(ctx, offset, pageSize)
}

func (s *inventoryService) UpdatePart(ctx context.Context, id uuid.UUID, req CreateSparePartRequest) (*models.SparePart, error) {
	p, err := s.repo.FindPartByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("inventoryService.UpdatePart find: %w", err)
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Category != "" {
		p.Category = req.Category
	}
	if req.Unit != "" {
		p.Unit = req.Unit
	}
	if req.MinStockLevel >= 0 {
		p.MinStockLevel = req.MinStockLevel
	}
	if req.UnitCost >= 0 {
		p.UnitCost = req.UnitCost
	}
	if err := s.repo.UpdatePart(ctx, p); err != nil {
		return nil, fmt.Errorf("inventoryService.UpdatePart: %w", err)
	}
	return p, nil
}

func (s *inventoryService) StockIn(ctx context.Context, partID uuid.UUID, userID uuid.UUID, req StockAdjustRequest) error {
	movement := &models.StockMovement{
		ID:          uuid.New(),
		SparePartID: partID,
		Type:        models.StockIn,
		Quantity:    req.Quantity,
		Date:        time.Now(),
		Notes:       req.Notes,
		CreatedBy:   userID,
	}
	if req.ReferenceID != "" {
		rid, _ := uuid.Parse(req.ReferenceID)
		movement.ReferenceID = &rid
		movement.ReferenceType = req.ReferenceType
	}
	if err := s.repo.CreateMovement(ctx, movement); err != nil {
		return fmt.Errorf("inventoryService.StockIn movement: %w", err)
	}
	if err := s.repo.AdjustStock(ctx, partID, req.Quantity); err != nil {
		return fmt.Errorf("inventoryService.StockIn adjust: %w", err)
	}
	return nil
}

func (s *inventoryService) StockOut(ctx context.Context, partID uuid.UUID, userID uuid.UUID, req StockAdjustRequest) error {
	part, err := s.repo.FindPartByID(ctx, partID)
	if err != nil {
		return fmt.Errorf("inventoryService.StockOut find: %w", err)
	}
	if part.CurrentStock < req.Quantity {
		return fmt.Errorf("insufficient stock: available %.2f, requested %.2f", part.CurrentStock, req.Quantity)
	}
	movement := &models.StockMovement{
		ID:          uuid.New(),
		SparePartID: partID,
		Type:        models.StockOut,
		Quantity:    req.Quantity,
		Date:        time.Now(),
		Notes:       req.Notes,
		CreatedBy:   userID,
	}
	if req.ReferenceID != "" {
		rid, _ := uuid.Parse(req.ReferenceID)
		movement.ReferenceID = &rid
		movement.ReferenceType = req.ReferenceType
	}
	if err := s.repo.CreateMovement(ctx, movement); err != nil {
		return fmt.Errorf("inventoryService.StockOut movement: %w", err)
	}
	if err := s.repo.AdjustStock(ctx, partID, -req.Quantity); err != nil {
		return fmt.Errorf("inventoryService.StockOut adjust: %w", err)
	}
	return nil
}

func (s *inventoryService) ListLowStock(ctx context.Context) ([]models.SparePart, error) {
	return s.repo.ListLowStock(ctx)
}

func (s *inventoryService) ListMovements(ctx context.Context, partID *uuid.UUID, page, pageSize int) ([]models.StockMovement, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListMovements(ctx, partID, offset, pageSize)
}
