package app

import (
	pb "github.com/daniarmas/api-example/pkg"
	"github.com/daniarmas/api-example/usecase"
)

type ItemServer struct {
	pb.UnimplementedItemServiceServer
	itemService usecase.ItemService
}

type AuthenticationServer struct {
	pb.UnimplementedAuthenticationServiceServer
	authenticationService usecase.AuthenticationService
}

func NewItemServer(
	itemService usecase.ItemService,
) *ItemServer {
	return &ItemServer{
		itemService: itemService,
	}
}

func NewAuthenticationServer(
	authenticationService usecase.AuthenticationService,
) *AuthenticationServer {
	return &AuthenticationServer{
		authenticationService: authenticationService,
	}
}
