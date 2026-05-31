package handlers

import (
	"time"

	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/services"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/agrifleet/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FinanceHandler struct {
	svc services.FinanceService
}

func NewFinanceHandler(svc services.FinanceService) *FinanceHandler {
	return &FinanceHandler{svc: svc}
}

// --- Expenses ---

func (h *FinanceHandler) ListExpenses(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var category *models.ExpenseCategory
	if cat := c.Query("category"); cat != "" {
		ec := models.ExpenseCategory(cat)
		category = &ec
	}
	var machineID *uuid.UUID
	if mid := c.Query("machine_id"); mid != "" {
		id, _ := uuid.Parse(mid)
		machineID = &id
	}
	var from, to *time.Time
	if f := c.Query("from"); f != "" {
		t, _ := time.Parse("2006-01-02", f)
		from = &t
	}
	if t := c.Query("to"); t != "" {
		parsed, _ := time.Parse("2006-01-02", t)
		to = &parsed
	}
	expenses, total, err := h.svc.ListExpenses(c.Context(), category, machineID, from, to, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch expenses")
	}
	return response.SuccessWithMeta(c, "Expenses retrieved", expenses, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *FinanceHandler) CreateExpense(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	var req services.CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	expense, err := h.svc.CreateExpense(c.Context(), claims.UserID, req)
	if err != nil {
		return response.InternalError(c, "Failed to create expense")
	}
	return response.Created(c, "Expense created", expense)
}

func (h *FinanceHandler) GetExpense(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}
	expense, err := h.svc.GetExpense(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Expense not found")
	}
	return response.Success(c, "Expense retrieved", expense)
}

func (h *FinanceHandler) UpdateExpense(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}
	var req services.CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	expense, err := h.svc.UpdateExpense(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update expense")
	}
	return response.Success(c, "Expense updated", expense)
}

func (h *FinanceHandler) DeleteExpense(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid expense ID")
	}
	if err := h.svc.DeleteExpense(c.Context(), id); err != nil {
		return response.InternalError(c, "Failed to delete expense")
	}
	return response.Success(c, "Expense deleted", nil)
}

// --- Revenues ---

func (h *FinanceHandler) ListRevenues(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var from, to *time.Time
	if f := c.Query("from"); f != "" {
		t, _ := time.Parse("2006-01-02", f)
		from = &t
	}
	if t := c.Query("to"); t != "" {
		parsed, _ := time.Parse("2006-01-02", t)
		to = &parsed
	}
	revenues, total, err := h.svc.ListRevenues(c.Context(), from, to, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch revenues")
	}
	return response.SuccessWithMeta(c, "Revenues retrieved", revenues, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *FinanceHandler) CreateRevenue(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	var req services.CreateRevenueRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	revenue, err := h.svc.CreateRevenue(c.Context(), claims.UserID, req)
	if err != nil {
		return response.InternalError(c, "Failed to create revenue")
	}
	return response.Created(c, "Revenue created", revenue)
}

// --- Summary ---

func (h *FinanceHandler) GetPLSummary(c *fiber.Ctx) error {
	from, to := defaultDateRange(c)
	summary, err := h.svc.GetPLSummary(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch P&L summary")
	}
	return response.Success(c, "P&L summary retrieved", summary)
}

func (h *FinanceHandler) GetMachineProfitability(c *fiber.Ctx) error {
	from, to := defaultDateRange(c)
	data, err := h.svc.GetMachineProfitability(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch machine profitability")
	}
	return response.Success(c, "Machine profitability retrieved", data)
}

func defaultDateRange(c *fiber.Ctx) (time.Time, time.Time) {
	from := time.Now().AddDate(0, -1, 0)
	to := time.Now()
	if f := c.Query("from"); f != "" {
		from, _ = time.Parse("2006-01-02", f)
	}
	if t := c.Query("to"); t != "" {
		to, _ = time.Parse("2006-01-02", t)
	}
	return from, to
}
