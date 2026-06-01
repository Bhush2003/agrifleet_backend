package handlers

import (
	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditHandler struct{ db *gorm.DB }

func NewAuditHandler(db *gorm.DB) *AuditHandler { return &AuditHandler{db: db} }

func (h *AuditHandler) ListLogs(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var logs []models.AuditLog
	var total int64
	q := h.db.WithContext(c.Context()).Model(&models.AuditLog{})
	if m := c.Query("module"); m != "" {
		q = q.Where("module = ?", m)
	}
	if u := c.Query("user_id"); u != "" {
		q = q.Where("user_id = ?", u)
	}
	q.Count(&total)
	q.Preload("User").Offset(p.Offset).Limit(p.PageSize).Order("created_at DESC").Find(&logs)
	return response.SuccessWithMeta(c, "Audit logs retrieved", logs, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *AuditHandler) ListNotifications(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	userIDStr := c.Query("user_id")
	var notifications []models.Notification
	var total int64
	q := h.db.WithContext(c.Context()).Model(&models.Notification{})
	if userIDStr != "" {
		q = q.Where("user_id = ?", userIDStr)
	}
	if unread := c.Query("unread"); unread == "true" {
		q = q.Where("is_read = false")
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PageSize).Order("created_at DESC").Find(&notifications)
	return response.SuccessWithMeta(c, "Notifications retrieved", notifications, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *AuditHandler) MarkNotificationRead(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid notification ID")
	}
	h.db.WithContext(c.Context()).Model(&models.Notification{}).
		Where("id = ?", id).Update("is_read", true)
	return response.Success(c, "Notification marked as read", nil)
}

func (h *AuditHandler) ListTasks(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var tasks []models.Task
	var total int64
	q := h.db.WithContext(c.Context()).Model(&models.Task{})
	if status := c.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if assignee := c.Query("assigned_to"); assignee != "" {
		q = q.Where("assigned_to = ?", assignee)
	}
	q.Count(&total)
	q.Preload("Assignee").Preload("Creator").
		Offset(p.Offset).Limit(p.PageSize).Order("created_at DESC").Find(&tasks)
	return response.SuccessWithMeta(c, "Tasks retrieved", tasks, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *AuditHandler) CreateTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	task.ID = uuid.New()
	task.Status = "open"
	if err := h.db.WithContext(c.Context()).Create(&task).Error; err != nil {
		return response.InternalError(c, "Failed to create task")
	}
	return response.Created(c, "Task created", task)
}

func (h *AuditHandler) UpdateTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid task ID")
	}
	var task models.Task
	if err := h.db.WithContext(c.Context()).First(&task, "id = ?", id).Error; err != nil {
		return response.NotFound(c, "Task not found")
	}
	if err := c.BodyParser(&task); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	h.db.WithContext(c.Context()).Save(&task)
	return response.Success(c, "Task updated", task)
}

func (h *AuditHandler) ListMonthLocks(c *fiber.Ctx) error {
	var locks []models.MonthLock
	h.db.WithContext(c.Context()).Order("year DESC, month DESC").Find(&locks)
	return response.Success(c, "Month locks retrieved", locks)
}

func (h *AuditHandler) ToggleMonthLock(c *fiber.Ctx) error {
	var body struct {
		Month  int    `json:"month"`
		Year   int    `json:"year"`
		Module string `json:"module"`
		Lock   bool   `json:"lock"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	var lock models.MonthLock
	h.db.WithContext(c.Context()).FirstOrCreate(&lock,
		models.MonthLock{Month: body.Month, Year: body.Year, Module: body.Module})
	lock.IsLocked = body.Lock
	h.db.WithContext(c.Context()).Save(&lock)
	msg := "Month unlocked"
	if body.Lock {
		msg = "Month locked"
	}
	return response.Success(c, msg, lock)
}

func (h *AuditHandler) ListRecycleBin(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var items []models.RecycleBin
	var total int64
	q := h.db.WithContext(c.Context()).Model(&models.RecycleBin{}).
		Where("is_restored = false")
	if et := c.Query("entity_type"); et != "" {
		q = q.Where("entity_type = ?", et)
	}
	q.Count(&total)
	q.Preload("Deleter").Offset(p.Offset).Limit(p.PageSize).
		Order("deleted_at DESC").Find(&items)
	return response.SuccessWithMeta(c, "Recycle bin retrieved", items, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *AuditHandler) GetSettings(c *fiber.Ctx) error {
	var settings []models.SystemSetting
	q := h.db.WithContext(c.Context())
	if group := c.Query("group"); group != "" {
		q = q.Where("\"group\" = ?", group)
	}
	q.Find(&settings)
	return response.Success(c, "Settings retrieved", settings)
}

func (h *AuditHandler) UpdateSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	var body struct {
		Value string `json:"value"`
	}
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	var setting models.SystemSetting
	h.db.WithContext(c.Context()).FirstOrCreate(&setting, models.SystemSetting{Key: key})
	setting.Value = body.Value
	h.db.WithContext(c.Context()).Save(&setting)
	return response.Success(c, "Setting updated", setting)
}

func (h *AuditHandler) ListDocuments(c *fiber.Ctx) error {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")
	var docs []models.Document
	q := h.db.WithContext(c.Context())
	if entityType != "" {
		q = q.Where("entity_type = ?", entityType)
	}
	if entityID != "" {
		q = q.Where("entity_id = ?", entityID)
	}
	q.Order("expiry_date ASC").Find(&docs)
	return response.Success(c, "Documents retrieved", docs)
}

func (h *AuditHandler) CreateDocument(c *fiber.Ctx) error {
	var doc models.Document
	if err := c.BodyParser(&doc); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	doc.ID = uuid.New()
	if err := h.db.WithContext(c.Context()).Create(&doc).Error; err != nil {
		return response.InternalError(c, "Failed to create document")
	}
	return response.Created(c, "Document created", doc)
}

func (h *AuditHandler) GetLedger(c *fiber.Ctx) error {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")
	if entityType == "" || entityID == "" {
		return response.BadRequest(c, "entity_type and entity_id are required")
	}
	p := pagination.FromContext(c)
	var entries []models.Ledger
	var total int64
	q := h.db.WithContext(c.Context()).Model(&models.Ledger{}).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID)
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PageSize).Order("date DESC").Find(&entries)
	return response.SuccessWithMeta(c, "Ledger retrieved", entries, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}
