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
	farmerH *handlers.FarmerHandler,
	factoryH *handlers.FactoryHandler,
	approvalH *handlers.ApprovalHandler,
	auditH *handlers.AuditHandler,
) {
	api := app.Group("/api/v1")

	// ── Auth (public) ──────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Post("/login", authH.Login)
	auth.Post("/refresh", authH.Refresh)

	// ── Protected ─────────────────────────────────────────────────────────────
	p := api.Use(middleware.AuthMiddleware(jwtSecret))

	p.Post("/auth/logout", authH.Logout)
	p.Get("/auth/me", authH.Me)

	// ── Machines ───────────────────────────────────────────────────────────────
	machines := p.Group("/machines")
	machines.Get("/", machineH.List)
	machines.Post("/",
		middleware.RequirePermission(models.PermVehicleCreate), machineH.Create)
	machines.Get("/:id", machineH.GetByID)
	machines.Put("/:id",
		middleware.RequirePermission(models.PermVehicleUpdate), machineH.Update)
	machines.Delete("/:id",
		middleware.RequirePermission(models.PermVehicleDelete), machineH.Delete)
	machines.Get("/:id/logs", machineH.ListLogs)
	machines.Post("/:id/logs", machineH.AddLog)

	// ── Workforce ──────────────────────────────────────────────────────────────
	drivers := p.Group("/drivers")
	drivers.Get("/", workforceH.ListDrivers)
	drivers.Post("/",
		middleware.RequirePermission(models.PermEmployeeCreate), workforceH.CreateDriver)
	drivers.Get("/:id", workforceH.GetDriver)
	drivers.Put("/:id",
		middleware.RequirePermission(models.PermEmployeeUpdate), workforceH.UpdateDriver)

	workers := p.Group("/workers")
	workers.Get("/", workforceH.ListWorkers)
	workers.Post("/",
		middleware.RequirePermission(models.PermEmployeeCreate), workforceH.CreateWorker)
	workers.Get("/:id", workforceH.GetWorker)
	workers.Put("/:id",
		middleware.RequirePermission(models.PermEmployeeUpdate), workforceH.UpdateWorker)

	p.Post("/attendance",
		middleware.RequirePermission(models.PermAttendanceMark), workforceH.MarkAttendance)
	p.Get("/attendance", workforceH.ListAttendance)
	p.Get("/payroll",
		middleware.RequirePermission(models.PermPayrollView), workforceH.GetPayroll)

	// ── Fuel ───────────────────────────────────────────────────────────────────
	fuel := p.Group("/fuel")
	fuel.Get("/", fuelH.List)
	fuel.Post("/",
		middleware.RequirePermission(models.PermDieselCreate), fuelH.Create)
	fuel.Get("/analytics", fuelH.GetAnalytics)
	fuel.Get("/:id", fuelH.GetByID)
	fuel.Put("/:id",
		middleware.RequirePermission(models.PermDieselCreate), fuelH.Update)
	fuel.Delete("/:id",
		middleware.RequirePermission(models.PermDieselApprove), fuelH.Delete)

	// ── Maintenance ────────────────────────────────────────────────────────────
	maint := p.Group("/maintenance")
	maint.Get("/", maintenanceH.List)
	maint.Post("/", maintenanceH.Create)
	maint.Get("/upcoming", maintenanceH.ListUpcoming)
	maint.Get("/overdue", maintenanceH.ListOverdue)
	maint.Get("/:id", maintenanceH.GetByID)
	maint.Put("/:id", maintenanceH.Update)

	// ── Inventory ──────────────────────────────────────────────────────────────
	parts := p.Group("/spare-parts")
	parts.Get("/", inventoryH.ListParts)
	parts.Post("/",
		middleware.RequirePermission(models.PermExpenseCreate), inventoryH.CreatePart)
	parts.Get("/low-stock", inventoryH.ListLowStock)
	parts.Get("/:id", inventoryH.GetPart)
	parts.Put("/:id",
		middleware.RequirePermission(models.PermExpenseCreate), inventoryH.UpdatePart)
	parts.Post("/:id/stock-in", inventoryH.StockIn)
	parts.Post("/:id/stock-out", inventoryH.StockOut)
	p.Get("/stock-movements", inventoryH.ListMovements)

	// ── Harvesting ─────────────────────────────────────────────────────────────
	harvest := p.Group("/harvesting/jobs")
	harvest.Get("/", harvestingH.ListJobs)
	harvest.Post("/",
		middleware.RequirePermission(models.PermIncomeCreate), harvestingH.CreateJob)
	harvest.Get("/:id", harvestingH.GetJob)
	harvest.Put("/:id",
		middleware.RequirePermission(models.PermIncomeCreate), harvestingH.UpdateJob)
	harvest.Post("/:id/logs", harvestingH.AddLog)
	harvest.Get("/:id/logs", harvestingH.ListLogs)
	harvest.Get("/:id/summary", harvestingH.GetJobSummary)

	// ── Transport ──────────────────────────────────────────────────────────────
	transport := p.Group("/transport")
	transport.Get("/trips", transportH.List)
	transport.Post("/trips", transportH.Create)
	transport.Get("/summary", transportH.GetSummary)
	transport.Get("/trips/:id", transportH.GetByID)
	transport.Put("/trips/:id", transportH.Update)

	// ── Finance ────────────────────────────────────────────────────────────────
	expenses := p.Group("/expenses")
	expenses.Get("/", financeH.ListExpenses)
	expenses.Post("/",
		middleware.RequirePermission(models.PermExpenseCreate), financeH.CreateExpense)
	expenses.Get("/:id", financeH.GetExpense)
	expenses.Put("/:id",
		middleware.RequirePermission(models.PermExpenseCreate), financeH.UpdateExpense)
	expenses.Delete("/:id",
		middleware.RequirePermission(models.PermExpenseApprove), financeH.DeleteExpense)

	revenues := p.Group("/revenues")
	revenues.Get("/", financeH.ListRevenues)
	revenues.Post("/",
		middleware.RequirePermission(models.PermIncomeCreate), financeH.CreateRevenue)

	p.Get("/finance/summary",
		middleware.RequirePermission(models.PermReportFinancial), financeH.GetPLSummary)
	p.Get("/finance/machine-profit",
		middleware.RequirePermission(models.PermReportFinancial), financeH.GetMachineProfitability)

	// ── Reports ────────────────────────────────────────────────────────────────
	reports := p.Group("/reports")
	reports.Get("/dashboard", reportH.Dashboard)
	reports.Get("/machines",
		middleware.RequirePermission(models.PermReportView), reportH.MachineReport)
	reports.Get("/fuel",
		middleware.RequirePermission(models.PermReportView), reportH.FuelReport)
	reports.Get("/workforce",
		middleware.RequirePermission(models.PermReportView), reportH.WorkforceReport)
	reports.Get("/projects",
		middleware.RequirePermission(models.PermReportView), reportH.ProjectReport)

	// ── Farmers ────────────────────────────────────────────────────────────────
	farmers := p.Group("/farmers")
	farmers.Get("/", farmerH.List)
	farmers.Post("/",
		middleware.RequirePermission(models.PermIncomeCreate), farmerH.Create)
	farmers.Get("/:id", farmerH.GetByID)
	farmers.Put("/:id",
		middleware.RequirePermission(models.PermIncomeCreate), farmerH.Update)
	farmers.Delete("/:id",
		middleware.RequirePermission(models.PermIncomeApprove), farmerH.Delete)
	farmers.Get("/:id/ledger", farmerH.GetLedger)

	// ── Factories ──────────────────────────────────────────────────────────────
	factories := p.Group("/factories")
	factories.Get("/", factoryH.List)
	factories.Post("/",
		middleware.RequirePermission(models.PermIncomeCreate), factoryH.Create)
	factories.Put("/:id",
		middleware.RequirePermission(models.PermIncomeCreate), factoryH.Update)
	factories.Delete("/:id",
		middleware.RequirePermission(models.PermIncomeApprove), factoryH.Delete)
	factories.Get("/:id/weight-slips", factoryH.ListWeightSlips)
	factories.Post("/:id/weight-slips", factoryH.CreateWeightSlip)

	// ── Approvals ──────────────────────────────────────────────────────────────
	approvals := p.Group("/approvals")
	approvals.Get("/", approvalH.List)
	approvals.Post("/", approvalH.Submit)
	approvals.Put("/:id/review",
		middleware.RequirePermission(models.PermApprovalManage), approvalH.Review)

	// ── Audit & Admin ──────────────────────────────────────────────────────────
	p.Get("/audit/logs",
		middleware.RequirePermission(models.PermAuditView), auditH.ListLogs)

	// Notifications
	p.Get("/notifications", auditH.ListNotifications)
	p.Put("/notifications/:id/read", auditH.MarkNotificationRead)

	// Tasks
	tasks := p.Group("/tasks")
	tasks.Get("/", auditH.ListTasks)
	tasks.Post("/", auditH.CreateTask)
	tasks.Put("/:id", auditH.UpdateTask)

	// Month locks
	p.Get("/month-locks", auditH.ListMonthLocks)
	p.Post("/month-locks",
		middleware.RequirePermission(models.PermMonthLock), auditH.ToggleMonthLock)

	// Recycle bin
	p.Get("/recycle-bin",
		middleware.RequirePermission(models.PermRecycleBin), auditH.ListRecycleBin)

	// Documents
	p.Get("/documents", auditH.ListDocuments)
	p.Post("/documents", auditH.CreateDocument)

	// Ledger
	p.Get("/ledger", auditH.GetLedger)

	// Settings (admin only)
	settings := p.Group("/settings",
		middleware.RequirePermission(models.PermSettingsManage))
	settings.Get("/", auditH.GetSettings)
	settings.Put("/:key", auditH.UpdateSetting)
}
