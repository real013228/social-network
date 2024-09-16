package resolvers

import (
	"context"
	"github.com/real013228/social-network/internal/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService userService
}

func NewResolver(userService userService) *Resolver {
	return &Resolver{userService: userService}
}

type userService interface {
	CreateUser(ctx context.Context, user model.CreateUserInput) (string, error)
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
}

type postService interface {
	CreatePost(ctx context.Context, post model.CreatePostInput) (string, error)
	GetPosts(ctx context.Context) ([]model.Post, error)
	GetPostsByFilter(ctx context.Context, filter model.PostsFilter) ([]model.Post, error)
}
