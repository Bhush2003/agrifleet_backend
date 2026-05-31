package handlers

import (
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

type MachineHandler struct {
	svc services.MachineService
}

func NewMachineHandler(svc services.MachineService) *MachineHandler {
	return &MachineHandler{svc: svc}
}

func (h *MachineHandler) List(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	filters := map[string]interface{}{}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if machineType := c.Query("type"); machineType != "" {
		filters["type"] = machineType
	}

	machines, total, err := h.svc.List(c.Context(), filters, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch machines")
	}
	return response.SuccessWithMeta(c, "Machines retrieved", machines, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *MachineHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	var req services.CreateMachineRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}

	machine, err := h.svc.Create(c.Context(), claims.UserID, req)
	if err != nil {
		return response.InternalError(c, "Failed to create machine")
	}
	return response.Created(c, "Machine created", machine)
}

func (h *MachineHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid machine ID")
	}
	machine, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Machine not found")
	}
	return response.Success(c, "Machine retrieved", machine)
}

func (h *MachineHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid machine ID")
	}
	var req services.UpdateMachineRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	machine, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update machine")
	}
	return response.Success(c, "Machine updated", machine)
}

func (h *MachineHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid machine ID")
	}
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return response.InternalError(c, "Failed to delete machine")
	}
	return response.Success(c, "Machine deleted", nil)
}

func (h *MachineHandler) AddLog(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	machineID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid machine ID")
	}
	var req services.AddMachineLogRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	log, err := h.svc.AddLog(c.Context(), machineID, claims.UserID, req)
	if err != nil {
		return response.InternalError(c, "Failed to add log")
	}
	return response.Created(c, "Log added", log)
}

func (h *MachineHandler) ListLogs(c *fiber.Ctx) error {
	machineID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid machine ID")
	}
	p := pagination.FromContext(c)
	logs, total, err := h.svc.ListLogs(c.Context(), machineID, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch logs")
	}
	return response.SuccessWithMeta(c, "Logs retrieved", logs, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

// Ensure models import is used (for status constants in query)
var _ = models.MachineStatusActive
