package initializers

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	externalURL := "postgres://nge_db_user:gvDhoCCHguRmZ7EhNpXtuH69sKlv6sIZ@dpg-cos9h7nsc6pc73e2rcc0-a.singapore-postgres.render.com/nge_db"
	DB, err = gorm.Open(postgres.Open(externalURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("Connected Successfully to the Database")
}

func GetDB() *gorm.DB {
	return DB
}
