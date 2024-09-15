package post_service

import (
	"context"
	"github.com/real013228/social-network/internal/model"
)

type postStorage interface {
	CreatePost(ctx context.Context, post model.CreatePostInput) (string, error)
	GetPosts(ctx context.Context) ([]model.Post, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error)
	GetPostByID(ctx context.Context, postID string) (model.Post, error)
	GetPostWithAllowedComments(ctx context.Context) ([]model.Post, error)
}

type PostService struct {
	storage postStorage
}

func NewPostService(storage postStorage) *PostService {
	return &PostService{storage: storage}
}
