package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinanceRepository interface {
	CreateExpense(ctx context.Context, e *models.Expense) error
	FindExpenseByID(ctx context.Context, id uuid.UUID) (*models.Expense, error)
	ListExpenses(ctx context.Context, category *models.ExpenseCategory, machineID *uuid.UUID, from, to *time.Time, offset, limit int) ([]models.Expense, int64, error)
	UpdateExpense(ctx context.Context, e *models.Expense) error
	DeleteExpense(ctx context.Context, id uuid.UUID) error

	CreateRevenue(ctx context.Context, r *models.Revenue) error
	FindRevenueByID(ctx context.Context, id uuid.UUID) (*models.Revenue, error)
	ListRevenues(ctx context.Context, from, to *time.Time, offset, limit int) ([]models.Revenue, int64, error)

	GetPLSummary(ctx context.Context, from, to time.Time) (*PLSummary, error)
	GetMachineProfitability(ctx context.Context, from, to time.Time) ([]MachineProfitability, error)
}

// PLSummary holds profit & loss totals for a period.
type PLSummary struct {
	TotalRevenue float64            `json:"total_revenue"`
	TotalExpense float64            `json:"total_expense"`
	NetProfit    float64            `json:"net_profit"`
	ByCategory   []ExpenseByCategory `json:"by_category"`
}

type ExpenseByCategory struct {
	Category models.ExpenseCategory `json:"category"`
	Amount   float64                `json:"amount"`
}

// MachineProfitability shows revenue vs expense per machine.
type MachineProfitability struct {
	MachineID   uuid.UUID `json:"machine_id"`
	MachineName string    `json:"machine_name"`
	Revenue     float64   `json:"revenue"`
	Expense     float64   `json:"expense"`
	Profit      float64   `json:"profit"`
}

type financeRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) FinanceRepository {
	return &financeRepository{db: db}
}

func (r *financeRepository) CreateExpense(ctx context.Context, e *models.Expense) error {
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return fmt.Errorf("financeRepository.CreateExpense: %w", err)
	}
	return nil
}

func (r *financeRepository) FindExpenseByID(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	var e models.Expense
	if err := r.db.WithContext(ctx).Preload("Machine").First(&e, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("financeRepository.FindExpenseByID: %w", err)
	}
	return &e, nil
}

func (r *financeRepository) ListExpenses(ctx context.Context, category *models.ExpenseCategory, machineID *uuid.UUID, from, to *time.Time, offset, limit int) ([]models.Expense, int64, error) {
	var expenses []models.Expense
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Expense{})
	if category != nil {
		q = q.Where("category = ?", *category)
	}
	if machineID != nil {
		q = q.Where("machine_id = ?", *machineID)
	}
	if from != nil {
		q = q.Where("date >= ?", *from)
	}
	if to != nil {
		q = q.Where("date <= ?", *to)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("financeRepository.ListExpenses count: %w", err)
	}
	if err := q.Preload("Machine").Offset(offset).Limit(limit).Order("date DESC").Find(&expenses).Error; err != nil {
		return nil, 0, fmt.Errorf("financeRepository.ListExpenses: %w", err)
	}
	return expenses, total, nil
}

func (r *financeRepository) UpdateExpense(ctx context.Context, e *models.Expense) error {
	if err := r.db.WithContext(ctx).Save(e).Error; err != nil {
		return fmt.Errorf("financeRepository.UpdateExpense: %w", err)
	}
	return nil
}

func (r *financeRepository) DeleteExpense(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Expense{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("financeRepository.DeleteExpense: %w", err)
	}
	return nil
}

func (r *financeRepository) CreateRevenue(ctx context.Context, rev *models.Revenue) error {
	if err := r.db.WithContext(ctx).Create(rev).Error; err != nil {
		return fmt.Errorf("financeRepository.CreateRevenue: %w", err)
	}
	return nil
}

func (r *financeRepository) FindRevenueByID(ctx context.Context, id uuid.UUID) (*models.Revenue, error) {
	var rev models.Revenue
	if err := r.db.WithContext(ctx).First(&rev, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("financeRepository.FindRevenueByID: %w", err)
	}
	return &rev, nil
}

func (r *financeRepository) ListRevenues(ctx context.Context, from, to *time.Time, offset, limit int) ([]models.Revenue, int64, error) {
	var revenues []models.Revenue
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Revenue{})
	if from != nil {
		q = q.Where("date >= ?", *from)
	}
	if to != nil {
		q = q.Where("date <= ?", *to)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("financeRepository.ListRevenues count: %w", err)
	}
	if err := q.Offset(offset).Limit(limit).Order("date DESC").Find(&revenues).Error; err != nil {
		return nil, 0, fmt.Errorf("financeRepository.ListRevenues: %w", err)
	}
	return revenues, total, nil
}

func (r *financeRepository) GetPLSummary(ctx context.Context, from, to time.Time) (*PLSummary, error) {
	var totalRevenue, totalExpense float64

	r.db.WithContext(ctx).Model(&models.Revenue{}).
		Where("date BETWEEN ? AND ?", from, to).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)

	r.db.WithContext(ctx).Model(&models.Expense{}).
		Where("date BETWEEN ? AND ?", from, to).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalExpense)

	var byCategory []ExpenseByCategory
	r.db.WithContext(ctx).Model(&models.Expense{}).
		Select("category, COALESCE(SUM(amount), 0) AS amount").
		Where("date BETWEEN ? AND ?", from, to).
		Group("category").
		Scan(&byCategory)

	return &PLSummary{
		TotalRevenue: totalRevenue,
		TotalExpense: totalExpense,
		NetProfit:    totalRevenue - totalExpense,
		ByCategory:   byCategory,
	}, nil
}

func (r *financeRepository) GetMachineProfitability(ctx context.Context, from, to time.Time) ([]MachineProfitability, error) {
	var results []MachineProfitability
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			m.id AS machine_id,
			m.name AS machine_name,
			COALESCE(rev.total, 0) AS revenue,
			COALESCE(exp.total, 0) AS expense,
			COALESCE(rev.total, 0) - COALESCE(exp.total, 0) AS profit
		FROM machines m
		LEFT JOIN (
			SELECT machine_id, SUM(amount) AS total
			FROM expenses
			WHERE date BETWEEN ? AND ? AND deleted_at IS NULL
			GROUP BY machine_id
		) exp ON exp.machine_id = m.id
		LEFT JOIN (
			SELECT source_id AS machine_id, SUM(amount) AS total
			FROM revenues
			WHERE date BETWEEN ? AND ? AND deleted_at IS NULL AND source_type = 'harvesting'
			GROUP BY source_id
		) rev ON rev.machine_id = m.id
		WHERE m.deleted_at IS NULL
		ORDER BY profit DESC
	`, from, to, from, to).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("financeRepository.GetMachineProfitability: %w", err)
	}
	return results, nil
}
