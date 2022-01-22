package dto

import "github.com/daniarmas/api-example/models"

type SignInResponse struct {
	RefreshToken       string
	AuthorizationToken string
	User               models.User
}
