package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/daniarmas/api-example/graph/generated"
	"github.com/daniarmas/api-example/graph/model"
)

func (r *queryResolver) Items(ctx context.Context, id string) ([]*model.Item, error) {
	items, err := r.ItemService.ListItem()
	if err != nil {
		return nil, err
	}
	itemsResponse := make([]*model.Item, 0, len(*items))
	for _, item := range *items {
		id := item.ID.String()
		cursor := int(item.Cursor)
		createTime := item.CreateTime.String()
		updateTime := item.UpdateTime.String()
		itemsResponse = append(itemsResponse, &model.Item{
			ID:                       &id,
			Name:                     &item.Name,
			Price:                    &item.Price,
			HighQualityPhoto:         &item.HighQualityPhoto,
			HighQualityPhotoBlurHash: &item.HighQualityPhotoBlurHash,
			LowQualityPhoto:          &item.LowQualityPhoto,
			LowQualityPhotoBlurHash:  &item.LowQualityPhotoBlurHash,
			Thumbnail:                &item.Thumbnail,
			ThumbnailBlurHash:        &item.ThumbnailBlurHash,
			CreateTime:               &createTime,
			UpdateTime:               &updateTime,
			Cursor:                   &cursor,
		})
	}
	return itemsResponse, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
