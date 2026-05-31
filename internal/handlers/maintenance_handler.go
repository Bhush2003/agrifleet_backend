package handlers

import (
	"github.com/agrifleet/backend/internal/services"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/agrifleet/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MaintenanceHandler struct {
	svc services.MaintenanceService
}

func NewMaintenanceHandler(svc services.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{svc: svc}
}

func (h *MaintenanceHandler) List(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var machineID *uuid.UUID
	if mid := c.Query("machine_id"); mid != "" {
		id, _ := uuid.Parse(mid)
		machineID = &id
	}
	records, total, err := h.svc.List(c.Context(), machineID, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch maintenance records")
	}
	return response.SuccessWithMeta(c, "Maintenance records retrieved", records, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *MaintenanceHandler) Create(c *fiber.Ctx) error {
	var req services.CreateMaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	record, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return response.InternalError(c, "Failed to create maintenance record")
	}
	return response.Created(c, "Maintenance record created", record)
}

func (h *MaintenanceHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid maintenance ID")
	}
	record, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Maintenance record not found")
	}
	return response.Success(c, "Maintenance record retrieved", record)
}

func (h *MaintenanceHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid maintenance ID")
	}
	var req services.CreateMaintenanceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	record, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update maintenance record")
	}
	return response.Success(c, "Maintenance record updated", record)
}

func (h *MaintenanceHandler) ListUpcoming(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30)
	records, err := h.svc.ListUpcoming(c.Context(), days)
	if err != nil {
		return response.InternalError(c, "Failed to fetch upcoming maintenance")
	}
	return response.Success(c, "Upcoming maintenance retrieved", records)
}

func (h *MaintenanceHandler) ListOverdue(c *fiber.Ctx) error {
	records, err := h.svc.ListOverdue(c.Context())
	if err != nil {
		return response.InternalError(c, "Failed to fetch overdue maintenance")
	}
	return response.Success(c, "Overdue maintenance retrieved", records)
}
