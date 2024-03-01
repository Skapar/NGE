package config

import (
	"log"

	initializers "github.com/Skapar/NGE/pkg/nge/database/initializers"
	// "github.com/Skapar/NGE/pkg/nge/models"
)

func Connect() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}


// func main() {
// 	initializers.DB.AutoMigrate(&models.User{})
// 	fmt.Println("? Migration complete")
// }

