package seeds

import (
	"github.com/daniarmas/api-example/models"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, email string, password string) error {
	return db.Create(&models.User{Email: email, Password: password}).Error
}
