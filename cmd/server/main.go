package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/agrifleet/backend/internal/config"
	"github.com/agrifleet/backend/internal/database"
	"github.com/agrifleet/backend/internal/handlers"
	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/agrifleet/backend/internal/router"
	"github.com/agrifleet/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()

	// ── Database ───────────────────────────────────────────────────────────────
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// ── Redis (optional) ───────────────────────────────────────────────────────
	_, err = database.NewRedis(cfg)
	if err != nil {
		log.Warn().Err(err).Msg("Redis unavailable — caching disabled")
	}

	// ── Repositories ───────────────────────────────────────────────────────────
	userRepo        := repositories.NewUserRepository(db)
	machineRepo     := repositories.NewMachineRepository(db)
	workforceRepo   := repositories.NewWorkforceRepository(db)
	fuelRepo        := repositories.NewFuelRepository(db)
	maintenanceRepo := repositories.NewMaintenanceRepository(db)
	inventoryRepo   := repositories.NewInventoryRepository(db)
	harvestingRepo  := repositories.NewHarvestingRepository(db)
	transportRepo   := repositories.NewTransportRepository(db)
	financeRepo     := repositories.NewFinanceRepository(db)

	// ── Services ───────────────────────────────────────────────────────────────
	authSvc        := services.NewAuthService(userRepo, cfg)
	machineSvc     := services.NewMachineService(machineRepo)
	workforceSvc   := services.NewWorkforceService(workforceRepo)
	fuelSvc        := services.NewFuelService(fuelRepo)
	maintenanceSvc := services.NewMaintenanceService(maintenanceRepo)
	inventorySvc   := services.NewInventoryService(inventoryRepo)
	harvestingSvc  := services.NewHarvestingService(harvestingRepo)
	transportSvc   := services.NewTransportService(transportRepo)
	financeSvc     := services.NewFinanceService(financeRepo)

	// ── Handlers ───────────────────────────────────────────────────────────────
	authH        := handlers.NewAuthHandler(authSvc)
	machineH     := handlers.NewMachineHandler(machineSvc)
	workforceH   := handlers.NewWorkforceHandler(workforceSvc)
	fuelH        := handlers.NewFuelHandler(fuelSvc)
	maintenanceH := handlers.NewMaintenanceHandler(maintenanceSvc)
	inventoryH   := handlers.NewInventoryHandler(inventorySvc)
	harvestingH  := handlers.NewHarvestingHandler(harvestingSvc)
	transportH   := handlers.NewTransportHandler(transportSvc)
	financeH     := handlers.NewFinanceHandler(financeSvc)
	reportH      := handlers.NewReportHandler(machineSvc, fuelSvc, workforceSvc, financeSvc, transportSvc, harvestingSvc)

	// New enterprise handlers (use DB directly for simplicity)
	farmerH   := handlers.NewFarmerHandler(db)
	factoryH  := handlers.NewFactoryHandler(db)
	approvalH := handlers.NewApprovalHandler(db)
	auditH    := handlers.NewAuditHandler(db)

	// ── Fiber App ──────────────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: globalErrorHandler,
	})

	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(middleware.CORS(cfg.CORS.Origins))
	app.Use(middleware.Logger())

	// Rate limit login
	app.Use("/api/v1/auth/login", limiter.New(limiter.Config{Max: 10, Expiration: 60}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": cfg.App.Name})
	})

	// Register all routes
	router.Setup(
		app, cfg.JWT.Secret,
		authH, machineH, workforceH, fuelH,
		maintenanceH, inventoryH, harvestingH,
		transportH, financeH, reportH,
		farmerH, factoryH, approvalH, auditH,
	)

	// ── Start ──────────────────────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.App.Port
	}
	if port == "" {
		port = "8080"
	}

	go func() {
		log.Info().Str("port", port).Msg("AgriFleet API starting")
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down...")
	_ = app.Shutdown()
}

func globalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{"success": false, "message": err.Error()})
}
