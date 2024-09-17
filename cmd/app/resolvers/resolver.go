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
	GetUsers(ctx context.Context, filter model.UsersFilter) ([]model.User, error)
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
	GetNotifications(ctx context.Context, filter model.UsersFilter) ([]string, error)
}

type postService interface {
	CreatePost(ctx context.Context, post model.CreatePostInput) (string, error)
	GetPosts(ctx context.Context, filter model.PostsFilter) ([]model.Post, error)
	GetPostsByFilter(ctx context.Context, filter model.PostsFilter) ([]model.Post, error)
	Subscribe(ctx context.Context, subscribeInput model.SubscribeInput) (model.SubscribePayload, error)
}

type commentService interface {
	CreateComment(ctx context.Context, comment model.CreateCommentInput) (string, error)
	GetReplies(ctx context.Context, commentID string) ([]model.Comment, error)
	GetComments(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error)
	GetCommentsByAuthorID(ctx context.Context, authorID string) ([]model.Comment, error)
}
