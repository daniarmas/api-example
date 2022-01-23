package seeds

import (
	"github.com/daniarmas/api-example/models"
	"gorm.io/gorm"
)

func CreateItem(db *gorm.DB, name string, price float64, highQualityPhoto string, highQualityPhotoBlurHash string, lowQualityPhoto string, lowQualityPhotoBlurHash string, thumbnail string, thumbnailBlurHash string) error {
	return db.Create(&models.Item{Name: name, Price: price, HighQualityPhoto: highQualityPhoto, HighQualityPhotoBlurHash: highQualityPhotoBlurHash, LowQualityPhoto: lowQualityPhoto, LowQualityPhotoBlurHash: lowQualityPhotoBlurHash, Thumbnail: thumbnail, ThumbnailBlurHash: thumbnailBlurHash}).Error
}
