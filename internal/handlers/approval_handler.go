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

type ApprovalHandler struct {
	db *gorm.DB
}

func NewApprovalHandler(db *gorm.DB) *ApprovalHandler {
	return &ApprovalHandler{db: db}
}

func (h *ApprovalHandler) List(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	role := models.UserRole(claims.Role)

	var requests []models.ApprovalRequest
	var total int64
	q := h.db.WithContext(c.Context()).Model(&models.ApprovalRequest{})

	// Managers see their own submissions; admins see all
	if role != models.RoleAdmin && role != models.RoleOwner {
		q = q.Where("submitted_by = ?", claims.UserID)
	}
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if atype := c.Query("type"); atype != "" {
		q = q.Where("type = ?", atype)
	}

	q.Count(&total)
	q.Preload("Submitter").Preload("Reviewer").
		Offset(p.Offset).Limit(p.PageSize).
		Order("created_at DESC").Find(&requests)

	return response.SuccessWithMeta(c, "Approval requests retrieved", requests, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *ApprovalHandler) Submit(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	var req models.ApprovalRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	req.ID = uuid.New()
	req.SubmittedBy = claims.UserID
	req.Status = models.ApprovalPending
	if err := h.db.WithContext(c.Context()).Create(&req).Error; err != nil {
		return response.InternalError(c, "Failed to submit approval request")
	}
	// Create notification for admins
	h.notifyAdmins(c, req)
	return response.Created(c, "Approval request submitted", req)
}

func (h *ApprovalHandler) Review(c *fiber.Ctx) error {
	claims := c.Locals(middleware.LocalsUserKey).(*pkgjwt.Claims)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid request ID")
	}

	var body struct {
		Status     models.ApprovalStatus `json:"status"`
		ReviewNote string                `json:"review_note"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	now := time.Now()
	result := h.db.WithContext(c.Context()).Model(&models.ApprovalRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      body.Status,
			"review_note": body.ReviewNote,
			"reviewed_by": claims.UserID,
			"reviewed_at": now,
		})
	if result.Error != nil {
		return response.InternalError(c, "Failed to review request")
	}

	// Log audit
	h.logAudit(c, claims.UserID, "approval_review", "approval", &id,
		string(body.Status))

	return response.Success(c, "Approval request reviewed", nil)
}

func (h *ApprovalHandler) notifyAdmins(c *fiber.Ctx, req models.ApprovalRequest) {
	var admins []models.User
	h.db.WithContext(c.Context()).
		Where("role IN ? AND is_active = true", []string{"admin", "owner"}).
		Find(&admins)
	for _, admin := range admins {
		notif := models.Notification{
			ID:     uuid.New(),
			UserID: admin.ID,
			Title:  "New Approval Request",
			Body:   req.Title,
			Type:   "approval",
		}
		h.db.WithContext(c.Context()).Create(&notif)
	}
}

func (h *ApprovalHandler) logAudit(c *fiber.Ctx, userID uuid.UUID, action, module string, entityID *uuid.UUID, newVal string) {
	log := models.AuditLog{
		ID:         uuid.New(),
		UserID:     userID,
		Action:     action,
		Module:     module,
		EntityID:   entityID,
		EntityType: module,
		NewValue:   newVal,
		IPAddress:  c.IP(),
	}
	h.db.WithContext(c.Context()).Create(&log)
}
