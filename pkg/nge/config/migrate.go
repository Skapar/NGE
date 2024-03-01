package config

import (
	"fmt"
	"log"

	initializers "github.com/Skapar/NGE/pkg/nge/database/initializers"
	"github.com/Skapar/NGE/pkg/nge/models"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config) // ConnectDB now just initializes the DB variable

	// Use GetDB to retrieve the initialized *gorm.DB instance
	db := initializers.GetDB()
	if db == nil {
		log.Fatal("Database connection was not established")
	}

	fmt.Println("? Database connection established successfully")

	return db, nil
}

func Migrate() {
	initializers.DB.AutoMigrate(&models.Event{})
	fmt.Println("? Migration complete")
}
