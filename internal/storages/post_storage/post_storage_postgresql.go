package post_storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages"
	"github.com/real013228/social-network/internal/storages/user_storage"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type PostStoragePostgres struct {
	client      storages.Client
	userStorage user_storage.UserStoragePostgres
}

func (s *PostStoragePostgres) Subscribe(ctx context.Context, subscribeInput model.SubscribeInput) (string, error) {
	q := `
		INSERT INTO subscriptions (post_id, user_id)
		VALUES ($1, $2)
		RETURNING post_id
	`
	var postID string
	if err := s.client.QueryRow(ctx, q, subscribeInput.PostID, subscribeInput.UserID).Scan(&postID); err != nil {
		return "", err
	}
	return fmt.Sprintf("successfully subscribed to post %s", postID), nil
}

func (s *PostStoragePostgres) GetSubscribers(ctx context.Context, postID string) ([]model.User, error) {
	q := `
		SELECT user_id
		FROM subscriptions
		WHERE post_id = $1
	`
	rows, err := s.client.Query(ctx, q, postID)
	if err != nil {
		return nil, err
	}
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *PostStoragePostgres) CreatePost(ctx context.Context, post model.Post) (string, error) {
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

func (s *PostStoragePostgres) GetPosts(ctx context.Context, filter model.PostsFilter) ([]model.Post, error) {
	q := `
		SELECT id, title, description, author_id, comments_allowed FROM posts
    	LIMIT $1 OFFSET $2;
	`
	rows, err := s.client.Query(ctx, q, filter.PageLimit, filter.PageNumber*filter.PageLimit)
	if err != nil {
		return nil, err
	}
	posts := make([]model.Post, 0)
	for rows.Next() {
		var post model.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.AuthorID, &post.CommentsAllowed)
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

func (s *PostStoragePostgres) GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error) {
	q := `
		SELECT posts.id, title, description, author_id, comments_allowed
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
		err = rows.Scan(&post.ID, &post.Title, &post.Description, &post.AuthorID, &post.CommentsAllowed)
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

func (s *PostStoragePostgres) GetPostByID(ctx context.Context, postID string) (model.Post, error) {
	q := `
		SELECT id, title, description, author_id, comments_allowed
		FROM posts 
		WHERE id = $1;
	`
	var post model.Post
	if err := s.client.QueryRow(ctx, q, postID).Scan(&post.ID, &post.Title, &post.Description, &post.AuthorID, &post.CommentsAllowed); err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func (s *PostStoragePostgres) GetPostWithAllowedComments(ctx context.Context) ([]model.Post, error) {
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

func NewPostStoragePostgres(client storages.Client, userStorage user_storage.UserStoragePostgres) *PostStoragePostgres {
	return &PostStoragePostgres{client: client, userStorage: userStorage}
}
