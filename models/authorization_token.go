package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tabler interface {
	TableName() string
}

func (AuthorizationToken) TableName() string {
	return AuthorizationTokenTableName
}

const AuthorizationTokenTableName = "authorization_token"

type AuthorizationToken struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()"`
	RefreshTokenFk uuid.UUID      `gorm:"type:uuid;column:refresh_token_fk;unique;not null"`
	RefreshToken   RefreshToken   `gorm:"foreignKey:RefreshTokenFk"`
	UserFk         uuid.UUID      `gorm:"column:user_fk;not null"`
	User           User           `gorm:"foreignKey:UserFk"`
	CreateTime     time.Time      `gorm:"column:create_time;not null"`
	UpdateTime     time.Time      `gorm:"column:update_time;not null"`
	DeleteTime     gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *AuthorizationToken) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}
