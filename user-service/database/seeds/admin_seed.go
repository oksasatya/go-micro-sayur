package seeds

import (
	"errors"
	"log"
	"user-service/internal/core/domain/model"
	"user-service/utils/conv"

	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	// Use a session that does not auto-save associations to avoid recursion
	tx := db.Session(&gorm.Session{FullSaveAssociations: false})

	hash, err := conv.HashPassword("admin123")
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	// Fetch the role by name (only what's needed)
	var role model.Roles
	if err := tx.Where("name = ?", "Super Admin").First(&role).Error; err != nil {
		log.Fatalf("load role: %v", err)
	}

	const email = "superadmin@example.com"
	var admin model.User

	// Find by unique field only (no associations)
	if err := tx.Where("email = ?", email).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create user without associations first
			admin = model.User{
				Name:       "Super Admin",
				Email:      email,
				Password:   hash,
				IsVerified: true,
			}
			if err := tx.Omit("Roles.*").Create(&admin).Error; err != nil {
				log.Fatalf("create admin: %v", err)
			}
		} else {
			log.Fatalf("find admin: %v", err)
		}
	} else {
		// Update existing user scalar fields safely
		if err := tx.Model(&admin).Updates(map[string]any{
			"name":        "Super Admin",
			"password":    hash,
			"is_verified": true,
		}).Error; err != nil {
			log.Fatalf("update admin: %v", err)
		}
	}

	// Manage many-to-many via association API to avoid recursive traversals
	if err := tx.Model(&admin).Association("Roles").Replace([]model.Roles{role}); err != nil {
		// Don't crash the app on association errors; log and return so the service can still start
		log.Printf("assign role failed: %v", err)
		return
	}

	log.Printf("admin %s seeded successfully", admin.Name)
}
