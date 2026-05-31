package handlers

import (
	"time"

	"github.com/agrifleet/backend/internal/services"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/agrifleet/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WorkforceHandler struct {
	svc services.WorkforceService
}

func NewWorkforceHandler(svc services.WorkforceService) *WorkforceHandler {
	return &WorkforceHandler{svc: svc}
}

// --- Drivers ---

func (h *WorkforceHandler) ListDrivers(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	drivers, total, err := h.svc.ListDrivers(c.Context(), p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch drivers")
	}
	return response.SuccessWithMeta(c, "Drivers retrieved", drivers, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *WorkforceHandler) CreateDriver(c *fiber.Ctx) error {
	var req services.CreateDriverRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	driver, err := h.svc.CreateDriver(c.Context(), req)
	if err != nil {
		return response.InternalError(c, "Failed to create driver")
	}
	return response.Created(c, "Driver created", driver)
}

func (h *WorkforceHandler) GetDriver(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid driver ID")
	}
	driver, err := h.svc.GetDriver(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Driver not found")
	}
	return response.Success(c, "Driver retrieved", driver)
}

func (h *WorkforceHandler) UpdateDriver(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid driver ID")
	}
	var req services.CreateDriverRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	driver, err := h.svc.UpdateDriver(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update driver")
	}
	return response.Success(c, "Driver updated", driver)
}

// --- Workers ---

func (h *WorkforceHandler) ListWorkers(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	workers, total, err := h.svc.ListWorkers(c.Context(), p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch workers")
	}
	return response.SuccessWithMeta(c, "Workers retrieved", workers, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *WorkforceHandler) CreateWorker(c *fiber.Ctx) error {
	var req services.CreateWorkerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	worker, err := h.svc.CreateWorker(c.Context(), req)
	if err != nil {
		return response.InternalError(c, "Failed to create worker")
	}
	return response.Created(c, "Worker created", worker)
}

func (h *WorkforceHandler) GetWorker(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid worker ID")
	}
	worker, err := h.svc.GetWorker(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Worker not found")
	}
	return response.Success(c, "Worker retrieved", worker)
}

func (h *WorkforceHandler) UpdateWorker(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid worker ID")
	}
	var req services.CreateWorkerRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	worker, err := h.svc.UpdateWorker(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update worker")
	}
	return response.Success(c, "Worker updated", worker)
}

// --- Attendance ---

func (h *WorkforceHandler) MarkAttendance(c *fiber.Ctx) error {
	var req services.BulkAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	if err := h.svc.MarkAttendance(c.Context(), req); err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, "Attendance marked", nil)
}

func (h *WorkforceHandler) ListAttendance(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var workerID *uuid.UUID
	if wid := c.Query("worker_id"); wid != "" {
		id, _ := uuid.Parse(wid)
		workerID = &id
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
	records, total, err := h.svc.ListAttendance(c.Context(), workerID, from, to, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch attendance")
	}
	return response.SuccessWithMeta(c, "Attendance retrieved", records, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *WorkforceHandler) GetPayroll(c *fiber.Ctx) error {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		return response.BadRequest(c, "from and to date parameters are required")
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return response.BadRequest(c, "Invalid from date format (use YYYY-MM-DD)")
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return response.BadRequest(c, "Invalid to date format (use YYYY-MM-DD)")
	}
	summary, err := h.svc.GetPayroll(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to calculate payroll")
	}
	return response.Success(c, "Payroll calculated", summary)
}
