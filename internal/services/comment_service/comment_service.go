package comment_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages/post_storage"
	"github.com/real013228/social-network/internal/storages/user_storage"
)

type commentStorage interface {
	CreateComment(ctx context.Context, input model.Comment) (string, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]model.Comment, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]model.Comment, error)
}

type authorStorage interface {
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
}

type postStorage interface {
	GetPostByID(ctx context.Context, postID string) (model.Post, error)
}

type CommentService struct {
	storage       commentStorage
	authorStorage authorStorage
	postStorage   postStorage
}

func (c CommentService) CreateComment(ctx context.Context, comment model.CreateCommentInput) (string, error) {
	newID := uuid.New()
	authorID := comment.AuthorID
	_, err := c.authorStorage.GetUserByID(ctx, model.UsersFilter{UserID: &authorID})
	if err != nil {
		if !errors.Is(err, user_storage.ErrUserNotFound) {
			return "", fmt.Errorf("user_storage.GetAuthorByID: %w", err)
		}
		return "", user_storage.ErrUserNotFound
	}

	postID := comment.PostID
	_, err = c.postStorage.GetPostByID(ctx, postID)
	if err != nil {
		if !errors.Is(err, post_storage.ErrPostNotFound) {
			return "", fmt.Errorf("post_storage.GetPostByID: %w", err)
		}
		return "", post_storage.ErrPostNotFound
	}
	var comm = model.Comment{
		ID:       newID.String(),
		PostID:   comment.PostID,
		Text:     comment.Text,
		AuthorID: comment.AuthorID,
	}
	id, err := c.storage.CreateComment(ctx, comm)
	if err != nil {
		return "", fmt.Errorf("comment_storage.CreateComment: %w", err)
	}

	return id, nil
}

func (c CommentService) GetCommentsByPostID(ctx context.Context, postID string) ([]model.Comment, error) {
	if _, err := c.postStorage.GetPostByID(ctx, postID); err != nil {
		if !errors.Is(err, post_storage.ErrPostNotFound) {
			return nil, fmt.Errorf("post_storage.GetPostByID: %w", err)
		}
		return nil, post_storage.ErrPostNotFound
	}

	comms, err := c.storage.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	return comms, nil
}

func (c CommentService) GetCommentsByAuthorID(ctx context.Context, authorID string) ([]model.Comment, error) {
	if _, err := c.authorStorage.GetUserByID(ctx, model.UsersFilter{UserID: &authorID}); err != nil {
		if !errors.Is(err, user_storage.ErrUserNotFound) {
			return nil, fmt.Errorf("user_storage.GetUserByID: %w", err)
		}
		return nil, user_storage.ErrUserNotFound
	}
	comms, err := c.storage.GetCommentsByUserID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	return comms, nil
}

func NewCommentService(storage commentStorage, authorStorage authorStorage, postStorage postStorage) *CommentService {
	return &CommentService{storage: storage, authorStorage: authorStorage, postStorage: postStorage}
}
