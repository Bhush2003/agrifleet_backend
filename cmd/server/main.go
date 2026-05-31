package main

import (
	"fmt"
	"os"

	"github.com/agrifleet/backend/internal/config"
	"github.com/agrifleet/backend/internal/database"
	"github.com/agrifleet/backend/internal/handlers"
	"github.com/agrifleet/backend/internal/middleware"
	"github.com/agrifleet/backend/internal/repositories"
	"github.com/agrifleet/backend/internal/router"
	"github.com/agrifleet/backend/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	userRepo := repositories.NewUserRepository(db)
	machineRepo := repositories.NewMachineRepository(db)
	workforceRepo := repositories.NewWorkforceRepository(db)
	fuelRepo := repositories.NewFuelRepository(db)
	maintenanceRepo := repositories.NewMaintenanceRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	harvestingRepo := repositories.NewHarvestingRepository(db)
	transportRepo := repositories.NewTransportRepository(db)
	financeRepo := repositories.NewFinanceRepository(db)

	authSvc := services.NewAuthService(userRepo, cfg)
	machineSvc := services.NewMachineService(machineRepo)
	workforceSvc := services.NewWorkforceService(workforceRepo)
	fuelSvc := services.NewFuelService(fuelRepo)
	maintenanceSvc := services.NewMaintenanceService(maintenanceRepo)
	inventorySvc := services.NewInventoryService(inventoryRepo)
	harvestingSvc := services.NewHarvestingService(harvestingRepo)
	transportSvc := services.NewTransportService(transportRepo)
	financeSvc := services.NewFinanceService(financeRepo)

	authH := handlers.NewAuthHandler(authSvc)
	machineH := handlers.NewMachineHandler(machineSvc)
	workforceH := handlers.NewWorkforceHandler(workforceSvc)
	fuelH := handlers.NewFuelHandler(fuelSvc)
	maintenanceH := handlers.NewMaintenanceHandler(maintenanceSvc)
	inventoryH := handlers.NewInventoryHandler(inventorySvc)
	harvestingH := handlers.NewHarvestingHandler(harvestingSvc)
	transportH := handlers.NewTransportHandler(transportSvc)
	financeH := handlers.NewFinanceHandler(financeSvc)
	reportH := handlers.NewReportHandler(
		machineSvc, fuelSvc, workforceSvc, financeSvc, transportSvc, harvestingSvc,
	)

	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: errorHandler,
	})

	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS(cfg.CORS.Origins))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	router.Setup(
		app, cfg.JWT.Secret,
		authH, machineH, workforceH, fuelH, maintenanceH,
		inventoryH, harvestingH, transportH, financeH, reportH,
	)

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Info().Str("port", cfg.App.Port).Msg("starting server")
	if err := app.Listen(addr); err != nil {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{"error": err.Error()})
}
