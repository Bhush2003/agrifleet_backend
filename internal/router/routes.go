package router

import (
	"github.com/agrifleet/backend/internal/handlers"
	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

// Setup registers all API routes on the Fiber app.
func Setup(
	app *fiber.App,
	jwtSecret string,
	authH *handlers.AuthHandler,
	machineH *handlers.MachineHandler,
	workforceH *handlers.WorkforceHandler,
	fuelH *handlers.FuelHandler,
	maintenanceH *handlers.MaintenanceHandler,
	inventoryH *handlers.InventoryHandler,
	harvestingH *handlers.HarvestingHandler,
	transportH *handlers.TransportHandler,
	financeH *handlers.FinanceHandler,
	reportH *handlers.ReportHandler,
) {
	api := app.Group("/api/v1")

	// ── Auth (public) ──────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Post("/login", authH.Login)
	auth.Post("/refresh", authH.Refresh)

	// ── Protected routes ───────────────────────────────────────────────────────
	protected := api.Use(middleware.AuthMiddleware(jwtSecret))

	protected.Post("/auth/logout", authH.Logout)
	protected.Get("/auth/me", authH.Me)

	// ── Machines ───────────────────────────────────────────────────────────────
	machines := protected.Group("/machines")
	machines.Get("/", machineH.List)
	machines.Post("/", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), machineH.Create)
	machines.Get("/:id", machineH.GetByID)
	machines.Put("/:id", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), machineH.Update)
	machines.Delete("/:id", middleware.RequireRoles(models.RoleOwner), machineH.Delete)
	machines.Get("/:id/logs", machineH.ListLogs)
	machines.Post("/:id/logs", machineH.AddLog)

	// ── Workforce ──────────────────────────────────────────────────────────────
	drivers := protected.Group("/drivers")
	drivers.Get("/", workforceH.ListDrivers)
	drivers.Post("/", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), workforceH.CreateDriver)
	drivers.Get("/:id", workforceH.GetDriver)
	drivers.Put("/:id", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), workforceH.UpdateDriver)

	workers := protected.Group("/workers")
	workers.Get("/", workforceH.ListWorkers)
	workers.Post("/", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), workforceH.CreateWorker)
	workers.Get("/:id", workforceH.GetWorker)
	workers.Put("/:id", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), workforceH.UpdateWorker)

	protected.Post("/attendance", workforceH.MarkAttendance)
	protected.Get("/attendance", workforceH.ListAttendance)
	protected.Get("/payroll", middleware.RequireRoles(models.RoleOwner, models.RoleAccountant), workforceH.GetPayroll)

	// ── Fuel ───────────────────────────────────────────────────────────────────
	fuel := protected.Group("/fuel")
	fuel.Get("/", fuelH.List)
	fuel.Post("/", fuelH.Create)
	fuel.Get("/analytics", fuelH.GetAnalytics)
	fuel.Get("/:id", fuelH.GetByID)
	fuel.Put("/:id", fuelH.Update)
	fuel.Delete("/:id", middleware.RequireRoles(models.RoleOwner, models.RoleAccountant), fuelH.Delete)

	// ── Maintenance ────────────────────────────────────────────────────────────
	maintenance := protected.Group("/maintenance")
	maintenance.Get("/", maintenanceH.List)
	maintenance.Post("/", maintenanceH.Create)
	maintenance.Get("/upcoming", maintenanceH.ListUpcoming)
	maintenance.Get("/overdue", maintenanceH.ListOverdue)
	maintenance.Get("/:id", maintenanceH.GetByID)
	maintenance.Put("/:id", maintenanceH.Update)

	// ── Inventory ──────────────────────────────────────────────────────────────
	spareParts := protected.Group("/spare-parts")
	spareParts.Get("/", inventoryH.ListParts)
	spareParts.Post("/", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), inventoryH.CreatePart)
	spareParts.Get("/low-stock", inventoryH.ListLowStock)
	spareParts.Get("/:id", inventoryH.GetPart)
	spareParts.Put("/:id", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), inventoryH.UpdatePart)
	spareParts.Post("/:id/stock-in", inventoryH.StockIn)
	spareParts.Post("/:id/stock-out", inventoryH.StockOut)

	protected.Get("/stock-movements", inventoryH.ListMovements)

	// ── Harvesting ─────────────────────────────────────────────────────────────
	harvesting := protected.Group("/harvesting/jobs")
	harvesting.Get("/", harvestingH.ListJobs)
	harvesting.Post("/", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), harvestingH.CreateJob)
	harvesting.Get("/:id", harvestingH.GetJob)
	harvesting.Put("/:id", middleware.RequireRoles(models.RoleOwner, models.RoleSupervisor), harvestingH.UpdateJob)
	harvesting.Post("/:id/logs", harvestingH.AddLog)
	harvesting.Get("/:id/logs", harvestingH.ListLogs)
	harvesting.Get("/:id/summary", harvestingH.GetJobSummary)

	// ── Transport ──────────────────────────────────────────────────────────────
	transport := protected.Group("/transport")
	transport.Get("/trips", transportH.List)
	transport.Post("/trips", transportH.Create)
	transport.Get("/summary", transportH.GetSummary)
	transport.Get("/trips/:id", transportH.GetByID)
	transport.Put("/trips/:id", transportH.Update)

	// ── Finance ────────────────────────────────────────────────────────────────
	expenses := protected.Group("/expenses")
	expenses.Get("/", financeH.ListExpenses)
	expenses.Post("/", financeH.CreateExpense)
	expenses.Get("/:id", financeH.GetExpense)
	expenses.Put("/:id", financeH.UpdateExpense)
	expenses.Delete("/:id", middleware.RequireRoles(models.RoleOwner, models.RoleAccountant), financeH.DeleteExpense)

	revenues := protected.Group("/revenues")
	revenues.Get("/", financeH.ListRevenues)
	revenues.Post("/", financeH.CreateRevenue)

	protected.Get("/finance/summary", middleware.RequireRoles(models.RoleOwner, models.RoleAccountant), financeH.GetPLSummary)
	protected.Get("/finance/machine-profit", middleware.RequireRoles(models.RoleOwner, models.RoleAccountant), financeH.GetMachineProfitability)

	// ── Reports ────────────────────────────────────────────────────────────────
	reports := protected.Group("/reports")
	reports.Get("/dashboard", reportH.Dashboard)
	reports.Get("/machines", reportH.MachineReport)
	reports.Get("/fuel", reportH.FuelReport)
	reports.Get("/workforce", reportH.WorkforceReport)
	reports.Get("/projects", reportH.ProjectReport)
}
