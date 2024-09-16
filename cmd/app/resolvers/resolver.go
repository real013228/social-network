package resolvers

import (
	"context"
	"github.com/real013228/social-network/internal/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService    userService
	postService    postService
	commentService commentService
}

func NewResolver(userService userService, postService postService, commentService commentService) *Resolver {
	return &Resolver{userService: userService, postService: postService, commentService: commentService}
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

type commentService interface {
	CreateComment(ctx context.Context, comment model.CreateCommentInput) (string, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]model.Comment, error)
	GetCommentsByAuthorID(ctx context.Context, authorID string) ([]model.Comment, error)
}
