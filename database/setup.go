package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	database, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	// Don't need to migrate, because we are using sqlite
	// err = database.AutoMigrate(&User{})
	// if err != nil {
	// 	panic("Failed to migrate database!")
	// }

	DB = database
}
