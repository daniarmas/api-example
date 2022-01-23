package usecase

import (
	"github.com/daniarmas/api-example/models"
	"github.com/daniarmas/api-example/repository"
	"gorm.io/gorm"
)

type ItemService interface {
	GetItem(id string) (*models.Item, error)
	ListItem() (*[]models.Item, error)
}

type itemService struct {
	dao repository.DAO
}

func NewItemService(dao repository.DAO) ItemService {
	return &itemService{dao: dao}
}

func (i *itemService) ListItem() (*[]models.Item, error) {
	var items []models.Item
	var itemsErr error
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		items, itemsErr = i.dao.NewItemQuery().ListItem(tx, &models.Item{})
		if itemsErr != nil {
			return itemsErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func (i *itemService) GetItem(id string) (*models.Item, error) {
	item, err := i.dao.NewItemQuery().GetItem(id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
