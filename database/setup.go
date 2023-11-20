package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"seahorse.app/server/database/models"
)

var DB *gorm.DB

func SetupDatabase() {
	database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate database")
	}

	DB = database
}
