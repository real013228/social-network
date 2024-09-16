package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/real013228/social-network/cmd/app/resolvers"
	"github.com/real013228/social-network/graph"
	"github.com/real013228/social-network/internal/services/comment_service"
	"github.com/real013228/social-network/internal/services/post_service"
	"github.com/real013228/social-network/internal/services/user_service"
	"github.com/real013228/social-network/internal/storages"
	"github.com/real013228/social-network/internal/storages/comment_storage"
	"github.com/real013228/social-network/internal/storages/post_storage"
	"github.com/real013228/social-network/internal/storages/user_storage"
	"log"
	"net/http"
	"os"
)

// add comment
const defaultPort = "8080"

func initializePostgreSQLServer(cfg storages.StorageConfig) *handler.Server {
	postgreSQLClient, err := storages.NewClient(context.TODO(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	userStoragePostgres := user_storage.NewUserStoragePostgres(postgreSQLClient)
	userService := user_service.NewUserService(userStoragePostgres)
	postStoragePostgres := post_storage.NewPostStoragePostgres(postgreSQLClient, *userStoragePostgres)
	postService := post_service.NewPostService(postStoragePostgres, userStoragePostgres)
	commentStoragePostgres := comment_storage.NewCommentStoragePostgres(postgreSQLClient)
	commentService := comment_service.NewCommentService(commentStoragePostgres, userStoragePostgres, postStoragePostgres)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolvers.NewResolver(
		userService,
		postService,
		commentService,
	),
	}))

	return srv
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	cfg := storages.StorageConfig{
		Host:          os.Getenv("POSTGRES_HOST"),
		Port:          os.Getenv("POSTGRES_PORT"),
		Username:      os.Getenv("POSTGRES_USER"),
		Password:      os.Getenv("POSTGRES_PASSWORD"),
		DBName:        os.Getenv("POSTGRES_DBNAME"),
		RetryAttempts: 3,
	}

	srv := initializePostgreSQLServer(cfg)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
