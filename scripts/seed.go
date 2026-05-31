//go:build ignore

// Run with: go run scripts/seed.go
// Seeds the database with a default admin user.
package main

import (
	"fmt"
	"log"

	"github.com/agrifleet/backend/internal/config"
	"github.com/agrifleet/backend/internal/database"
	"github.com/agrifleet/backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Check if admin already exists
	var count int64
	db.Model(&models.User{}).Where("phone = ?", "+1234567890").Count(&count)
	if count > 0 {
		fmt.Println("✅ Admin user already exists — skipping seed")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Hash failed: %v", err)
	}

	users := []models.User{
		{
			ID:           uuid.New(),
			Name:         "Admin Owner",
			Email:        "admin@agrifleet.io",
			Phone:        "+1234567890",
			PasswordHash: string(hash),
			Role:         models.RoleOwner,
			IsActive:     true,
		},
		{
			ID:           uuid.New(),
			Name:         "Field Supervisor",
			Email:        "supervisor@agrifleet.io",
			Phone:        "+1234567891",
			PasswordHash: string(hash),
			Role:         models.RoleSupervisor,
			IsActive:     true,
		},
		{
			ID:           uuid.New(),
			Name:         "John Driver",
			Email:        "driver@agrifleet.io",
			Phone:        "+1234567892",
			PasswordHash: string(hash),
			Role:         models.RoleDriver,
			IsActive:     true,
		},
		{
			ID:           uuid.New(),
			Name:         "Accountant",
			Email:        "accountant@agrifleet.io",
			Phone:        "+1234567893",
			PasswordHash: string(hash),
			Role:         models.RoleAccountant,
			IsActive:     true,
		},
	}

	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			log.Printf("❌ Failed to create user %s: %v", u.Name, err)
		} else {
			fmt.Printf("✅ Created user: %s (%s) — phone: %s\n", u.Name, u.Role, u.Phone)
		}
	}

	fmt.Println("\n🌱 Seed complete!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Login credentials (all use password: password123)")
	fmt.Println("  Owner      → +1234567890")
	fmt.Println("  Supervisor → +1234567891")
	fmt.Println("  Driver     → +1234567892")
	fmt.Println("  Accountant → +1234567893")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
