package post_storage

import (
	"context"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages"
	"github.com/real013228/social-network/internal/storages/user_storage"
)

type PostStoragePostgres struct {
	client      storages.Client
	userStorage user_storage.UserStoragePostgres
}

func (s PostStoragePostgres) CreatePost(ctx context.Context, post model.Post) (string, error) {
	q := `
		INSERT INTO posts (id, title, description, author_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var id string
	if err := s.client.QueryRow(ctx, q, post.ID, post.Title, post.Description, post.AuthorID).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (s PostStoragePostgres) GetPosts(ctx context.Context) ([]model.Post, error) {
	q := `
		SELECT id, title, description, author_id FROM posts;
	`
	rows, err := s.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s PostStoragePostgres) GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error) {
	q := `
		SELECT id, title, description, author_id 
		FROM posts JOIN users ON users.id = posts.author_id
		WHERE users.id = $1;
	`
	var posts []model.Post
	rows, err := s.client.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s PostStoragePostgres) GetPostByID(ctx context.Context, postID string) (model.Post, error) {
	q := `
		SELECT id, title, description, author_id
		FROM posts 
		WHERE id = $1;
	`
	var post model.Post
	if err := s.client.QueryRow(ctx, q, postID).Scan(&post); err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func (s PostStoragePostgres) GetPostWithAllowedComments(ctx context.Context) ([]model.Post, error) {
	q := `
		SELECT id, title, description, author_id
		FROM posts
		WHERE comments_allowed = true;
	`
	var posts []model.Post
	rows, err := s.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
