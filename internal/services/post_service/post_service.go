package post_service

import "github.com/real013228/social-network/internal/model"

type postStorage interface {
	CreatePost(post model.Post) error
	GetPosts(filter model.PostsFilter) ([]model.Post, error)
}

type PostService struct {
	storage postStorage
}

func NewPostService(storage postStorage) *PostService {
	return &PostService{storage: storage}
}
