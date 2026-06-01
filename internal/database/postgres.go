package database

import (
	"fmt"

	"github.com/agrifleet/backend/internal/config"
	"github.com/agrifleet/backend/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgres creates and returns a GORM PostgreSQL connection.
func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User,
		cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode, cfg.DB.Timezone,
	)

	gormCfg := &gorm.Config{}
	if cfg.App.Env == "development" {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormCfg.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("database.NewPostgres: %w", err)
	}

	log.Info().Msg("PostgreSQL connected successfully")
	return db, nil
}

// AutoMigrate runs GORM auto-migration for all models.
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Machine{},
		&models.MachineLog{},
		&models.Driver{},
		&models.Worker{},
		&models.Attendance{},
		&models.FuelEntry{},
		&models.Maintenance{},
		&models.SparePart{},
		&models.StockMovement{},
		&models.HarvestingJob{},
		&models.HarvestingLog{},
		&models.TransportTrip{},
		&models.Expense{},
		&models.Revenue{},
		&models.Project{},
		&models.RefreshToken{},
		// New enterprise models
		&models.Farmer{},
		&models.FarmerLedger{},
		&models.Factory{},
		&models.WeightSlip{},
		&models.ApprovalRequest{},
		&models.AuditLog{},
		&models.MonthLock{},
		&models.RecycleBin{},
		&models.Notification{},
		&models.Task{},
		&models.Document{},
		&models.Ledger{},
		&models.SystemSetting{},
	)
	if err != nil {
		return fmt.Errorf("database.AutoMigrate: %w", err)
	}
	log.Info().Msg("Database migration completed")
	return nil
}
