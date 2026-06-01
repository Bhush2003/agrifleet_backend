package handlers

import (
	"time"

	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/models"
	pkgjwt "github.com/agrifleet/backend/pkg/jwt"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FarmerHandler struct {
	db *gorm.DB
}

func NewFarmerHandler(db *gorm.DB) *FarmerHandler {
	return &FarmerHandler{db: db}
}

func (h *FarmerHandler) List(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var farmers []models.Farmer
	var total int64
	q := h.db.WithContext(c.Context()).Model(&models.Farmer{})
	if search := c.Query("search"); search != "" {
		q = q.Where("name ILIKE ? OR mobile ILIKE ? OR village ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PageSize).Order("name ASC").Find(&farmers)
	return response.SuccessWithMeta(c, "Farmers retrieved", farmers, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *FarmerHandler) Create(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	var farmer models.Farmer
	if err := c.BodyParser(&farmer); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	farmer.ID = uuid.New()
	farmer.CreatedBy = claims.UserID
	if err := h.db.WithContext(c.Context()).Create(&farmer).Error; err != nil {
		return response.InternalError(c, "Failed to create farmer")
	}
	return response.Created(c, "Farmer created", farmer)
}

func (h *FarmerHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid farmer ID")
	}
	var farmer models.Farmer
	if err := h.db.WithContext(c.Context()).First(&farmer, "id = ?", id).Error; err != nil {
		return response.NotFound(c, "Farmer not found")
	}
	return response.Success(c, "Farmer retrieved", farmer)
}

func (h *FarmerHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid farmer ID")
	}
	var farmer models.Farmer
	if err := h.db.WithContext(c.Context()).First(&farmer, "id = ?", id).Error; err != nil {
		return response.NotFound(c, "Farmer not found")
	}
	if err := c.BodyParser(&farmer); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	h.db.WithContext(c.Context()).Save(&farmer)
	return response.Success(c, "Farmer updated", farmer)
}

func (h *FarmerHandler) Delete(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid farmer ID")
	}
	deletedBy := claims.UserID
	h.db.WithContext(c.Context()).Model(&models.Farmer{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"deleted_by": deletedBy})
	h.db.WithContext(c.Context()).Delete(&models.Farmer{}, "id = ?", id)
	return response.Success(c, "Farmer deleted", nil)
}

func (h *FarmerHandler) GetLedger(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid farmer ID")
	}
	p := pagination.FromContext(c)
	var entries []models.FarmerLedger
	var total int64
	h.db.WithContext(c.Context()).Model(&models.FarmerLedger{}).
		Where("farmer_id = ?", id).Count(&total)
	h.db.WithContext(c.Context()).
		Where("farmer_id = ?", id).
		Offset(p.Offset).Limit(p.PageSize).
		Order("date DESC").Find(&entries)
	return response.SuccessWithMeta(c, "Farmer ledger retrieved", entries, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

// Ensure time import is used
var _ = time.Now
