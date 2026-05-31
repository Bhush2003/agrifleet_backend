package handlers

import (
	"time"

	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/services"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/agrifleet/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FuelHandler struct {
	svc services.FuelService
}

func NewFuelHandler(svc services.FuelService) *FuelHandler {
	return &FuelHandler{svc: svc}
}

func (h *FuelHandler) List(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
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
	entries, total, err := h.svc.List(c.Context(), machineID, from, to, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch fuel entries")
	}
	return response.SuccessWithMeta(c, "Fuel entries retrieved", entries, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *FuelHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	var req services.CreateFuelEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	entry, err := h.svc.Create(c.Context(), claims.UserID, req)
	if err != nil {
		return response.InternalError(c, "Failed to create fuel entry")
	}
	return response.Created(c, "Fuel entry created", entry)
}

func (h *FuelHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid fuel entry ID")
	}
	entry, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Fuel entry not found")
	}
	return response.Success(c, "Fuel entry retrieved", entry)
}

func (h *FuelHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid fuel entry ID")
	}
	var req services.CreateFuelEntryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	entry, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update fuel entry")
	}
	return response.Success(c, "Fuel entry updated", entry)
}

func (h *FuelHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid fuel entry ID")
	}
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return response.InternalError(c, "Failed to delete fuel entry")
	}
	return response.Success(c, "Fuel entry deleted", nil)
}

func (h *FuelHandler) GetAnalytics(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	from := time.Now().AddDate(0, -1, 0)
	to := time.Now()
	if fromStr != "" {
		from, _ = time.Parse("2006-01-02", fromStr)
	}
	if toStr != "" {
		to, _ = time.Parse("2006-01-02", toStr)
	}
	analytics, err := h.svc.GetAnalytics(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch analytics")
	}
	return response.Success(c, "Fuel analytics retrieved", analytics)
}
