package post_storage

import (
	"context"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages"
)

type PostStoragePostgres struct {
	client storages.Client
}

func (s PostStoragePostgres) CreatePost(ctx context.Context, post model.CreatePostInput) (string, error) {
	q := `
		INSERT INTO posts (title, description, author_id)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	var id string
	if err := s.client.QueryRow(ctx, q, post.Title, post.Description, post.AuthorID).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
	//todo finish implementing posts, comments, and somehow make subscriptions
}

func (s PostStoragePostgres) GetPosts(ctx context.Context) ([]model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s PostStoragePostgres) GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s PostStoragePostgres) GetPostByID(ctx context.Context, postID string) (model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s PostStoragePostgres) GetPostWithAllowedComments(ctx context.Context) ([]model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func NewPostStoragePostgres(client *storages.Client) *PostStoragePostgres {
	return &PostStoragePostgres{client: client}
}
