package initializers

import (
	"backend/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}