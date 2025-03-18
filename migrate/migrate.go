package main

import (
	"FASMS/initializers"
	"FASMS/models"

	"log"
)

func init() {
	initializers.GetEnvs()
	initializers.ConnectDB()
}

func main() {
	if initializers.DB == nil {
		log.Fatal("Database connection is nil")
	}

	err := initializers.DB.AutoMigrate(&models.Applicants{})
	if err != nil {
		log.Fatal("Failed to migrate Applicants table:", err)
	}

	err = initializers.DB.AutoMigrate(&models.Households{})
	if err != nil {
		log.Fatal("Failed to migrate Households table:", err)
	}

	err = initializers.DB.AutoMigrate(&models.Applications{})
	if err != nil {
		log.Fatal("Failed to migrate Applications table:", err)
	}

	err = initializers.DB.AutoMigrate(&models.Schemes{})
	if err != nil {
		log.Fatal("Failed to migrate Schemes table:", err)
	}

	err = initializers.DB.AutoMigrate(&models.CriteriaGroup{})
	if err != nil {
		log.Fatal("Failed to migrate Criteria Group table:", err)
	}
	err = initializers.DB.AutoMigrate(&models.Criterias{})
	if err != nil {
		log.Fatal("Failed to migrate Criterias table:", err)
	}

	err = initializers.DB.AutoMigrate(&models.Benefits{})
	if err != nil {
		log.Fatal("Failed to migrate Benefits table:", err)
	}

}

//go mod migrate/migrate.go
