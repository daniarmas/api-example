package app

import (
	"context"

	"github.com/daniarmas/api-example/dto"
	pb "github.com/daniarmas/api-example/pkg"
)

func (m *ItemServer) ListItem(ctx context.Context, req *pb.ListItemRequest) (*pb.ListItemResponse, error) {
	items, err := m.itemService.ListItem(&dto.ListItemRequest{BusinessFk: req.BusinessFk, BusinessItemCategoryFk: req.ItemCategoryFk, NextPage: req.NextPage})
	if err != nil {
		return nil, err
	}
	itemsResponse := make([]*pb.Item, 0, len(*items))
	for _, item := range *items {
		itemsResponse = append(itemsResponse, &pb.Item{
			Id:                       item.ID.String(),
			Name:                     item.Name,
			Price:                    item.Price,
			HighQualityPhoto:         item.HighQualityPhoto,
			HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
			LowQualityPhoto:          item.LowQualityPhoto,
			LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
			Thumbnail:                item.Thumbnail,
			ThumbnailBlurHash:        item.ThumbnailBlurHash,
			CreateTime:               item.CreateTime.String(),
			UpdateTime:               item.UpdateTime.String(),
			Cursor:                   int32(item.Cursor),
		})
	}
	return &pb.ListItemResponse{Items: itemsResponse}, nil
}

func (m *ItemServer) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	item, err := m.itemService.GetItem(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetItemResponse{Item: &pb.Item{
		Id:                       item.ID.String(),
		Name:                     item.Name,
		Price:                    item.Price,
		HighQualityPhoto:         item.HighQualityPhoto,
		HighQualityPhotoBlurHash: item.HighQualityPhotoBlurHash,
		LowQualityPhoto:          item.LowQualityPhoto,
		LowQualityPhotoBlurHash:  item.LowQualityPhotoBlurHash,
		Thumbnail:                item.Thumbnail,
		ThumbnailBlurHash:        item.ThumbnailBlurHash,
		CreateTime:               item.CreateTime.String(),
		UpdateTime:               item.CreateTime.String(),
		Cursor:                   int32(item.Cursor),
	}}, nil
}
