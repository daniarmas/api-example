package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const RefreshTokenTableName = "refresh_token"

func (RefreshToken) TableName() string {
	return RefreshTokenTableName
}

type RefreshToken struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserFk     uuid.UUID      `gorm:"column:user_fk"`
	CreateTime time.Time      `gorm:"column:create_time"`
	UpdateTime time.Time      `gorm:"column:update_time"`
	DeleteTime gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	r.CreateTime = time.Now()
	r.UpdateTime = time.Now()
	return
}
