package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const UserTableName = "user"

func (User) TableName() string {
	return UserTableName
}

type User struct {
	ID                       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email                    string         `gorm:"column:email"`
	CreateTime               time.Time      `gorm:"column:create_time"`
	UpdateTime               time.Time      `gorm:"column:update_time"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()
	return
}
