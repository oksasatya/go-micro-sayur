package seeds

import (
	"gorm.io/gorm"
	"log"
	"user-service/internal/core/domain/model"
)

func SeedRole(db *gorm.DB) {
	// seed Role
	roles := []model.Roles{
		{
			Name: "Super Admin",
		},
		{
			Name: "Customer",
		},
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Roles{
			Name: role.Name,
		}).Error; err != nil {
			log.Fatalf("failed to seed role: "+err.Error(), err)
		} else {
			log.Printf("role %s seeded successfully", role.Name)
		}
	}
}
