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

type TransportHandler struct {
	svc services.TransportService
}

func NewTransportHandler(svc services.TransportService) *TransportHandler {
	return &TransportHandler{svc: svc}
}

func (h *TransportHandler) List(c *fiber.Ctx) error {
	p := pagination.FromContext(c)
	var jobID, truckID, driverID *uuid.UUID
	if v := c.Query("job_id"); v != "" {
		id, _ := uuid.Parse(v)
		jobID = &id
	}
	if v := c.Query("truck_id"); v != "" {
		id, _ := uuid.Parse(v)
		truckID = &id
	}
	if v := c.Query("driver_id"); v != "" {
		id, _ := uuid.Parse(v)
		driverID = &id
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
	trips, total, err := h.svc.List(c.Context(), jobID, truckID, driverID, from, to, p.Page, p.PageSize)
	if err != nil {
		return response.InternalError(c, "Failed to fetch trips")
	}
	return response.SuccessWithMeta(c, "Trips retrieved", trips, &response.Meta{
		Page: p.Page, PageSize: p.PageSize, Total: total,
		TotalPages: pagination.TotalPages(total, p.PageSize),
	})
}

func (h *TransportHandler) Create(c *fiber.Ctx) error {
	var req services.CreateTripRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	if errs := validator.Validate(req); errs != nil {
		return response.ValidationError(c, errs)
	}
	trip, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return response.InternalError(c, "Failed to create trip")
	}
	return response.Created(c, "Trip created", trip)
}

func (h *TransportHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid trip ID")
	}
	trip, err := h.svc.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "Trip not found")
	}
	return response.Success(c, "Trip retrieved", trip)
}

func (h *TransportHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid trip ID")
	}
	var req services.CreateTripRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}
	trip, err := h.svc.Update(c.Context(), id, req)
	if err != nil {
		return response.InternalError(c, "Failed to update trip")
	}
	return response.Success(c, "Trip updated", trip)
}

func (h *TransportHandler) GetSummary(c *fiber.Ctx) error {
	from := time.Now().AddDate(0, -1, 0)
	to := time.Now()
	if f := c.Query("from"); f != "" {
		from, _ = time.Parse("2006-01-02", f)
	}
	if t := c.Query("to"); t != "" {
		to, _ = time.Parse("2006-01-02", t)
	}
	summary, err := h.svc.GetSummary(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch transport summary")
	}
	return response.Success(c, "Transport summary retrieved", summary)
}
