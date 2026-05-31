package handlers

import (
	"github.com/agrifleet/backend/internal/models"
	"github.com/agrifleet/backend/internal/services"
	"github.com/agrifleet/backend/pkg/pagination"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/agrifleet/backend/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HarvestingHandler struct {
	svc services.HarvestingService
}

func NewHarvestingHandler(svc services.HarvestingService) *HarvestingHandler {
	return &HarvestingHandler{svc: svc}
}

func (h *HarvestingHandler) ListJobs(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var status *models.JobStatus
	if s := c.Query("status"); s != "" {
		js := models.JobStatus(s)
		status = &js
	}
	jobs, total, err := h.svc.ListJobs(c.Context(), status, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch harvesting jobs")
	}
	return response.SuccessWithMeta(c, "Harvesting jobs retrieved", jobs, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *HarvestingHandler) CreateJob(c *fiber.Ctx) error {
	var req services.CreateHarvestingJobRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	job, err := h.svc.CreateJob(c.Context(), req)
	if err != nil {
		return response.InternalError(c, "Failed to create harvesting job")
	}
	return response.Created(c, "Harvesting job created", job)
}

func (h *HarvestingHandler) GetJob(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid job ID")
	}
	job, err := h.svc.GetJob(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Harvesting job not found")
	}
	return response.Success(c, "Harvesting job retrieved", job)
}

func (h *HarvestingHandler) UpdateJob(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid job ID")
	}
	var req services.CreateHarvestingJobRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	job, err := h.svc.UpdateJob(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update harvesting job")
	}
	return response.Success(c, "Harvesting job updated", job)
}

func (h *HarvestingHandler) AddLog(c *fiber.Ctx) error {
	jobID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid job ID")
	}
	var req services.AddHarvestingLogRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	log, err := h.svc.AddLog(c.Context(), jobID, req)
	if err != nil {
		return response.InternalError(c, "Failed to add harvesting log")
	}
	return response.Created(c, "Harvesting log added", log)
}

func (h *HarvestingHandler) ListLogs(c *fiber.Ctx) error {
	jobID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid job ID")
	}
	p := pagination.FromContext(c)
	logs, total, err := h.svc.ListLogs(c.Context(), jobID, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch harvesting logs")
	}
	return response.SuccessWithMeta(c, "Harvesting logs retrieved", logs, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *HarvestingHandler) GetJobSummary(c *fiber.Ctx) error {
	jobID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid job ID")
	}
	summary, err := h.svc.GetJobSummary(c.Context(), jobID)
	if err != nil {
		return response.InternalError(c, "Failed to fetch job summary")
	}
	return response.Success(c, "Job summary retrieved", summary)
}
