package repositories

import (
	"context"
	"fmt"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	CreatePart(ctx context.Context, p *models.SparePart) error
	FindPartByID(ctx context.Context, id uuid.UUID) (*models.SparePart, error)
	ListParts(ctx context.Context, offset, limit int) ([]models.SparePart, int64, error)
	UpdatePart(ctx context.Context, p *models.SparePart) error
	ListLowStock(ctx context.Context) ([]models.SparePart, error)
	CreateMovement(ctx context.Context, m *models.StockMovement) error
	ListMovements(ctx context.Context, partID *uuid.UUID, offset, limit int) ([]models.StockMovement, int64, error)
	AdjustStock(ctx context.Context, partID uuid.UUID, delta float64) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) CreatePart(ctx context.Context, p *models.SparePart) error {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return fmt.Errorf("inventoryRepository.CreatePart: %w", err)
	}
	return nil
}

func (r *inventoryRepository) FindPartByID(ctx context.Context, id uuid.UUID) (*models.SparePart, error) {
	var p models.SparePart
	if err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("inventoryRepository.FindPartByID: %w", err)
	}
	return &p, nil
}

func (r *inventoryRepository) ListParts(ctx context.Context, offset, limit int) ([]models.SparePart, int64, error) {
	var parts []models.SparePart
	var total int64
	q := r.db.WithContext(ctx).Model(&models.SparePart{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("inventoryRepository.ListParts count: %w", err)
	}
	if err := q.Offset(offset).Limit(limit).Order("name ASC").Find(&parts).Error; err != nil {
		return nil, 0, fmt.Errorf("inventoryRepository.ListParts: %w", err)
	}
	return parts, total, nil
}

func (r *inventoryRepository) UpdatePart(ctx context.Context, p *models.SparePart) error {
	if err := r.db.WithContext(ctx).Save(p).Error; err != nil {
		return fmt.Errorf("inventoryRepository.UpdatePart: %w", err)
	}
	return nil
}

func (r *inventoryRepository) ListLowStock(ctx context.Context) ([]models.SparePart, error) {
	var parts []models.SparePart
	err := r.db.WithContext(ctx).
		Where("current_stock <= min_stock_level").
		Order("current_stock ASC").
		Find(&parts).Error
	if err != nil {
		return nil, fmt.Errorf("inventoryRepository.ListLowStock: %w", err)
	}
	return parts, nil
}

func (r *inventoryRepository) CreateMovement(ctx context.Context, m *models.StockMovement) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("inventoryRepository.CreateMovement: %w", err)
	}
	return nil
}

func (r *inventoryRepository) ListMovements(ctx context.Context, partID *uuid.UUID, offset, limit int) ([]models.StockMovement, int64, error) {
	var movements []models.StockMovement
	var total int64
	q := r.db.WithContext(ctx).Model(&models.StockMovement{})
	if partID != nil {
		q = q.Where("spare_part_id = ?", *partID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("inventoryRepository.ListMovements count: %w", err)
	}
	if err := q.Preload("SparePart").Offset(offset).Limit(limit).Order("date DESC").Find(&movements).Error; err != nil {
		return nil, 0, fmt.Errorf("inventoryRepository.ListMovements: %w", err)
	}
	return movements, total, nil
}

func (r *inventoryRepository) AdjustStock(ctx context.Context, partID uuid.UUID, delta float64) error {
	err := r.db.WithContext(ctx).
		Model(&models.SparePart{}).
		Where("id = ?", partID).
		UpdateColumn("current_stock", gorm.Expr("current_stock + ?", delta)).Error
	if err != nil {
		return fmt.Errorf("inventoryRepository.AdjustStock: %w", err)
	}
	return nil
}
