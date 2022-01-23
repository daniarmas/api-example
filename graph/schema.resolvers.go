package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/daniarmas/api-example/dto"
	"github.com/daniarmas/api-example/graph/generated"
	"github.com/daniarmas/api-example/graph/model"
)

func (r *mutationResolver) SignIn(ctx context.Context, email string, password string) (*model.SignInResponse, error) {
	result, err := r.AuthenticationService.SignIn(&dto.SignInRequest{Password: password, Email: email})
	if err != nil {
		switch err.Error() {
		case "user not found":
			return nil, errors.New("user not found")
		case "password incorrect":
			return nil, errors.New("credentials incorrect")
		default:
			return nil, errors.New("internal server error")
		}
	}
	id := result.User.ID.String()
	createTime := result.User.CreateTime.String()
	updateTime := result.User.UpdateTime.String()
	return &model.SignInResponse{RefreshToken: result.RefreshToken, AuthorizationToken: result.AuthorizationToken, User: &model.User{ID: &id, Email: &result.User.Email, CreateTime: &createTime, UpdateTime: &updateTime}}, nil
}

func (r *queryResolver) Items(ctx context.Context, id string) ([]*model.Item, error) {
	if id == "" {
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
	} else {
		item, err := r.ItemService.GetItem(id)
		if err != nil {
			return nil, err
		}
		id := item.ID.String()
		cursor := int(item.Cursor)
		createTime := item.CreateTime.String()
		updateTime := item.UpdateTime.String()
		itemsResponse := make([]*model.Item, 0, 1)
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
		return itemsResponse, nil
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
