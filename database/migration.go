package database

import (
	"log"
	"rakamin-final/models"
)

func Migrate() {
	Instance.AutoMigrate(&models.User{}, &models.Photo{})
	log.Println("Database Migration Completed!")
}
