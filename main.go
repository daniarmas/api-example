package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daniarmas/api-example/grpc"
	"github.com/daniarmas/api-example/graph"
	"github.com/daniarmas/api-example/graph/generated"
	"github.com/daniarmas/api-example/models"
	pb "github.com/daniarmas/api-example/pkg"
	"github.com/daniarmas/api-example/repository"
	"github.com/daniarmas/api-example/seeds"
	"github.com/daniarmas/api-example/usecase"
	"github.com/go-chi/chi"

	"google.golang.org/grpc"
)

func main() {
	// Load config file
	config, err := repository.NewConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// DB
	db, err := repository.NewDB(config)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return
	}
	// Register all services
	dao := repository.NewDAO(db, config)
	itemService := usecase.NewItemService(dao)
	authenticationService := usecase.NewAuthenticationService(dao)
	// Seed data
	result := db.Unscoped().Limit(1).Find(&models.User{})
	if result.RowsAffected == 0 {
		for _, seed := range seeds.All(&dao) {
			if err := seed.Run(db); err != nil {
				log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
			}
		}
	}
	// Starting graphQL server
	go func() {
		router := chi.NewRouter()
		srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{ItemService: itemService, AuthenticationService: authenticationService}}))
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
		router.Handle("/query", srv)
		graphqlPort := fmt.Sprintf(":%s", config.GraphqlApiPort)
		err := http.ListenAndServe(graphqlPort, router)
		if err != nil {
			panic(err)
		}
	}()
	// Starting gRPC server
	address := fmt.Sprintf("0.0.0.0:%s", config.GrpcApiPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	// Registring the services
	pb.RegisterItemServiceServer(grpcServer, app.NewItemServer(
		itemService,
	))
	pb.RegisterAuthenticationServiceServer(grpcServer, app.NewAuthenticationServer(
		authenticationService,
	))
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Server running at localhost:8081")
}
