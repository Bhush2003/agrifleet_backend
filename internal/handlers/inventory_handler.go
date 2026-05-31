package handlers

import (
	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/services"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/agrifleet/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	svc services.InventoryService
}

func NewInventoryHandler(svc services.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) ListParts(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	parts, total, err := h.svc.ListParts(c.Context(), p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch spare parts")
	}
	return response.SuccessWithMeta(c, "Spare parts retrieved", parts, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *InventoryHandler) CreatePart(c *fiber.Ctx) error {
	var req services.CreateSparePartRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	part, err := h.svc.CreatePart(c.Context(), req)
	if err != nil {
		return response.InternalError(c, "Failed to create spare part")
	}
	return response.Created(c, "Spare part created", part)
}

func (h *InventoryHandler) GetPart(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid spare part ID")
	}
	part, err := h.svc.GetPart(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Spare part not found")
	}
	return response.Success(c, "Spare part retrieved", part)
}

func (h *InventoryHandler) UpdatePart(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid spare part ID")
	}
	var req services.CreateSparePartRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	part, err := h.svc.UpdatePart(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update spare part")
	}
	return response.Success(c, "Spare part updated", part)
}

func (h *InventoryHandler) StockIn(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid spare part ID")
	}
	var req services.StockAdjustRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	if err := h.svc.StockIn(c.Context(), id, claims.UserID, req); err != nil {
		return response.InternalError(c, "Failed to record stock in")
	}
	return response.Success(c, "Stock in recorded", nil)
}

func (h *InventoryHandler) StockOut(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid spare part ID")
	}
	var req services.StockAdjustRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	if err := h.svc.StockOut(c.Context(), id, claims.UserID, req); err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Success(c, "Stock out recorded", nil)
}

func (h *InventoryHandler) ListLowStock(c *fiber.Ctx) error {
	parts, err := h.svc.ListLowStock(c.Context())
	if err != nil {
		return response.InternalError(c, "Failed to fetch low stock items")
	}
	return response.Success(c, "Low stock items retrieved", parts)
}

func (h *InventoryHandler) ListMovements(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var partID *uuid.UUID
	if pid := c.Query("spare_part_id"); pid != "" {
		id, _ := uuid.Parse(pid)
		partID = &id
	}
	movements, total, err := h.svc.ListMovements(c.Context(), partID, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch stock movements")
	}
	return response.SuccessWithMeta(c, "Stock movements retrieved", movements, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}
