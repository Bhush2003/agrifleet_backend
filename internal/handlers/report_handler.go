package handlers

import (
	"time"

	"github.com/agrifleet/backend/internal/services"
	"github.com/agrifleet/backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	machineSvc     services.MachineService
	fuelSvc        services.FuelService
	workforceSvc   services.WorkforceService
	financeSvc     services.FinanceService
	transportSvc   services.TransportService
	harvestingSvc  services.HarvestingService
}

func NewReportHandler(
	machineSvc services.MachineService,
	fuelSvc services.FuelService,
	workforceSvc services.WorkforceService,
	financeSvc services.FinanceService,
	transportSvc services.TransportService,
	harvestingSvc services.HarvestingService,
) *ReportHandler {
	return &ReportHandler{
		machineSvc:    machineSvc,
		fuelSvc:       fuelSvc,
		workforceSvc:  workforceSvc,
		financeSvc:    financeSvc,
		transportSvc:  transportSvc,
		harvestingSvc: harvestingSvc,
	}
}

// Dashboard returns KPI summary for the dashboard.
func (h *ReportHandler) Dashboard(c *fiber.Ctx) error {
	from := time.Now().AddDate(0, -1, 0)
	to := time.Now()

	pl, err := h.financeSvc.GetPLSummary(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch dashboard data")
	}

	fuelAnalytics, err := h.fuelSvc.GetAnalytics(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch fuel analytics")
	}

	machines, machineTotal, _ := h.machineSvc.List(c.Context(), map[string]interface{}{}, 1, 1000)
	activeMachines := 0
	for _, m := range machines {
		if m.Status == "active" {
			activeMachines++
		}
	}

	return response.Success(c, "Dashboard data retrieved", fiber.Map{
		"total_machines":  machineTotal,
		"active_machines": activeMachines,
		"pl_summary":      pl,
		"fuel_analytics":  fuelAnalytics,
		"period": fiber.Map{
			"from": from.Format("2006-01-02"),
			"to":   to.Format("2006-01-02"),
		},
	})
}

// MachineReport returns machine utilization and profitability.
func (h *ReportHandler) MachineReport(c *fiber.Ctx) error {
	from, to := defaultDateRange(c)
	profitability, err := h.financeSvc.GetMachineProfitability(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch machine report")
	}
	return response.Success(c, "Machine report retrieved", profitability)
}

// FuelReport returns fuel consumption analytics.
func (h *ReportHandler) FuelReport(c *fiber.Ctx) error {
	from, to := defaultDateRange(c)
	analytics, err := h.fuelSvc.GetAnalytics(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch fuel report")
	}
	return response.Success(c, "Fuel report retrieved", analytics)
}

// WorkforceReport returns attendance and payroll summary.
func (h *ReportHandler) WorkforceReport(c *fiber.Ctx) error {
	from, to := defaultDateRange(c)
	payroll, err := h.workforceSvc.GetPayroll(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch workforce report")
	}
	return response.Success(c, "Workforce report retrieved", payroll)
}

// ProjectReport returns project-wise P&L.
func (h *ReportHandler) ProjectReport(c *fiber.Ctx) error {
	from, to := defaultDateRange(c)
	pl, err := h.financeSvc.GetPLSummary(c.Context(), from, to)
	if err != nil {
		return response.InternalError(c, "Failed to fetch project report")
	}
	return response.Success(c, "Project report retrieved", pl)
}
