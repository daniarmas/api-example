package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const ItemTableName = "item"

func (Item) TableName() string {
	return ItemTableName
}

type Item struct {
	ID                       uuid.UUID 		`gorm:"type:uuid;default:uuid_generate_v4()"`
	Name                     string    		`gorm:"column:name;not null"`
	Price                    float64   		 `gorm:"column:price;not null"`
	HighQualityPhoto         string         `gorm:"column:high_quality_photo;not null"`
	HighQualityPhotoBlurHash string         `gorm:"column:high_quality_photo_blurhash;not null"`
	LowQualityPhoto          string         `gorm:"column:low_quality_photo;not null"`
	LowQualityPhotoBlurHash  string         `gorm:"column:low_quality_photo_blurhash;not null"`
	Thumbnail                string         `gorm:"column:thumbnail;not null"`
	ThumbnailBlurHash        string         `gorm:"column:thumbnail_blurhash;not null"`
	Cursor                   int64          `gorm:"column:cursor"`
	CreateTime               time.Time      `gorm:"column:create_time;not null"`
	UpdateTime               time.Time      `gorm:"column:update_time;not null"`
	DeleteTime               gorm.DeletedAt `gorm:"index;column:delete_time"`
}

func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now()
	i.UpdateTime = time.Now()
	return
}
