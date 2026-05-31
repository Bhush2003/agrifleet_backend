package services

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/google/uuid"
)

type CreateExpenseRequest struct {
	Category    string  `json:"category" validate:"required,oneof=fuel maintenance salary transport other"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Date        string  `json:"date" validate:"required"`
	MachineID   string  `json:"machine_id"`
	Description string  `json:"description"`
	PaymentMode string  `json:"payment_mode"`
}

type CreateRevenueRequest struct {
	SourceType  string  `json:"source_type" validate:"required,oneof=harvesting transport project other"`
	SourceID    string  `json:"source_id"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Date        string  `json:"date" validate:"required"`
	Description string  `json:"description"`
	PaymentMode string  `json:"payment_mode"`
}

type FinanceService interface {
	CreateExpense(ctx context.Context, userID uuid.UUID, req CreateExpenseRequest) (*models.Expense, error)
	GetExpense(ctx context.Context, id uuid.UUID) (*models.Expense, error)
	ListExpenses(ctx context.Context, category *models.ExpenseCategory, machineID *uuid.UUID, from, to *time.Time, page, pageSize int) ([]models.Expense, int64, error)
	UpdateExpense(ctx context.Context, id uuid.UUID, req CreateExpenseRequest) (*models.Expense, error)
	DeleteExpense(ctx context.Context, id uuid.UUID) error

	CreateRevenue(ctx context.Context, userID uuid.UUID, req CreateRevenueRequest) (*models.Revenue, error)
	GetRevenue(ctx context.Context, id uuid.UUID) (*models.Revenue, error)
	ListRevenues(ctx context.Context, from, to *time.Time, page, pageSize int) ([]models.Revenue, int64, error)

	GetPLSummary(ctx context.Context, from, to time.Time) (*repositories.PLSummary, error)
	GetMachineProfitability(ctx context.Context, from, to time.Time) ([]repositories.MachineProfitability, error)
}

type financeService struct {
	repo repositories.FinanceRepository
}

func NewFinanceService(repo repositories.FinanceRepository) FinanceService {
	return &financeService{repo: repo}
}

func (s *financeService) CreateExpense(ctx context.Context, userID uuid.UUID, req CreateExpenseRequest) (*models.Expense, error) {
	date, _ := time.Parse("2006-01-02", req.Date)
	e := &models.Expense{
		ID:          uuid.New(),
		Category:    models.ExpenseCategory(req.Category),
		Amount:      req.Amount,
		Date:        date,
		Description: req.Description,
		PaymentMode: models.PaymentMode(req.PaymentMode),
		CreatedBy:   userID,
	}
	if req.MachineID != "" {
		mid, _ := uuid.Parse(req.MachineID)
		e.MachineID = &mid
	}
	if err := s.repo.CreateExpense(ctx, e); err != nil {
		return nil, fmt.Errorf("financeService.CreateExpense: %w", err)
	}
	return e, nil
}

func (s *financeService) GetExpense(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	return s.repo.FindExpenseByID(ctx, id)
}

func (s *financeService) ListExpenses(ctx context.Context, category *models.ExpenseCategory, machineID *uuid.UUID, from, to *time.Time, page, pageSize int) ([]models.Expense, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListExpenses(ctx, category, machineID, from, to, offset, pageSize)
}

func (s *financeService) UpdateExpense(ctx context.Context, id uuid.UUID, req CreateExpenseRequest) (*models.Expense, error) {
	e, err := s.repo.FindExpenseByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("financeService.UpdateExpense find: %w", err)
	}
	if req.Amount > 0 {
		e.Amount = req.Amount
	}
	if req.Description != "" {
		e.Description = req.Description
	}
	if req.PaymentMode != "" {
		e.PaymentMode = models.PaymentMode(req.PaymentMode)
	}
	if err := s.repo.UpdateExpense(ctx, e); err != nil {
		return nil, fmt.Errorf("financeService.UpdateExpense: %w", err)
	}
	return e, nil
}

func (s *financeService) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteExpense(ctx, id)
}

func (s *financeService) CreateRevenue(ctx context.Context, userID uuid.UUID, req CreateRevenueRequest) (*models.Revenue, error) {
	date, _ := time.Parse("2006-01-02", req.Date)
	rev := &models.Revenue{
		ID:          uuid.New(),
		SourceType:  models.RevenueSource(req.SourceType),
		Amount:      req.Amount,
		Date:        date,
		Description: req.Description,
		PaymentMode: models.PaymentMode(req.PaymentMode),
		CreatedBy:   userID,
	}
	if req.SourceID != "" {
		sid, _ := uuid.Parse(req.SourceID)
		rev.SourceID = &sid
	}
	if err := s.repo.CreateRevenue(ctx, rev); err != nil {
		return nil, fmt.Errorf("financeService.CreateRevenue: %w", err)
	}
	return rev, nil
}

func (s *financeService) GetRevenue(ctx context.Context, id uuid.UUID) (*models.Revenue, error) {
	return s.repo.FindRevenueByID(ctx, id)
}

func (s *financeService) ListRevenues(ctx context.Context, from, to *time.Time, page, pageSize int) ([]models.Revenue, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.ListRevenues(ctx, from, to, offset, pageSize)
}

func (s *financeService) GetPLSummary(ctx context.Context, from, to time.Time) (*repositories.PLSummary, error) {
	return s.repo.GetPLSummary(ctx, from, to)
}

func (s *financeService) GetMachineProfitability(ctx context.Context, from, to time.Time) ([]repositories.MachineProfitability, error) {
	return s.repo.GetMachineProfitability(ctx, from, to)
}
