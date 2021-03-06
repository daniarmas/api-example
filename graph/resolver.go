package graph

import "github.com/daniarmas/api-example/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ItemService           usecase.ItemService
	AuthenticationService usecase.AuthenticationService
}
