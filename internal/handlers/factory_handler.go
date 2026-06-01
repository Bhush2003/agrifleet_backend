package handlers

import (
	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/models"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FactoryHandler struct{ db *gorm.DB }

func NewFactoryHandler(db *gorm.DB) *FactoryHandler { return &FactoryHandler{db: db} }

func (h *FactoryHandler) List(c *fiber.Ctx) error {
	var factories []models.Factory
	h.db.WithContext(c.Context()).Where("deleted_at IS NULL").Order("name ASC").Find(&factories)
	return response.Success(c, "Factories retrieved", factories)
}

func (h *FactoryHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	var factory models.Factory
	if err := c.BodyParser(&factory); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	factory.ID = uuid.New()
	factory.CreatedBy = claims.UserID
	if err := h.db.WithContext(c.Context()).Create(&factory).Error; err != nil {
		return response.InternalError(c, "Failed to create factory")
	}
	return response.Created(c, "Factory created", factory)
}

func (h *FactoryHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid factory ID")
	}
	var factory models.Factory
	if err := h.db.WithContext(c.Context()).First(&factory, "id = ?", id).Error; err != nil {
		return response.NotFound(c, "Factory not found")
	}
	if err := c.BodyParser(&factory); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	h.db.WithContext(c.Context()).Save(&factory)
	return response.Success(c, "Factory updated", factory)
}

func (h *FactoryHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid factory ID")
	}
	h.db.WithContext(c.Context()).Delete(&models.Factory{}, "id = ?", id)
	return response.Success(c, "Factory deleted", nil)
}

func (h *FactoryHandler) ListWeightSlips(c *fiber.Ctx) error {
	factoryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid factory ID")
	}
	p := pagination.FromContext(c)
	var slips []models.WeightSlip
	var total int64
	q := h.db.WithContext(c.Context()).Model(&models.WeightSlip{}).
		Where("factory_id = ?", factoryID)
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PageSize).Order("date DESC").Find(&slips)
	return response.SuccessWithMeta(c, "Weight slips retrieved", slips, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *FactoryHandler) CreateWeightSlip(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	var slip models.WeightSlip
	if err := c.BodyParser(&slip); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	slip.ID = uuid.New()
	slip.CreatedBy = claims.UserID
	slip.NetWeight = slip.GrossWeight - slip.TareWeight
	slip.TotalAmount = slip.NetWeight * slip.RatePerTon
	if err := h.db.WithContext(c.Context()).Create(&slip).Error; err != nil {
		return response.InternalError(c, "Failed to create weight slip")
	}
	return response.Created(c, "Weight slip created", slip)
}
