package comment_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages/comment_storage"
	"github.com/real013228/social-network/internal/storages/post_storage"
	"github.com/real013228/social-network/internal/storages/user_storage"
)

const (
	maxTextLength = 2000
)

var (
	ErrTextIsTooLarge = fmt.Errorf("text is too large, %d", maxTextLength)
)

type commentStorage interface {
	CreateComment(ctx context.Context, input model.Comment) (string, error)
	GetCommentsByPostID(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]model.Comment, error)
	GetCommentByID(ctx context.Context, commentID string) (model.Comment, error)
	GetReplies(ctx context.Context, commentID string) ([]model.Comment, error)
}

type authorStorage interface {
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
}

type postStorage interface {
	GetPostByID(ctx context.Context, postID string) (model.Post, error)
}

type postService interface {
	NotifyAll(ctx context.Context, payload model.NotificationPayload)
}

type CommentService struct {
	storage       commentStorage
	authorStorage authorStorage
	postStorage   postStorage
	postService   postService
}

func (c *CommentService) GetReplies(ctx context.Context, commentID string) ([]model.Comment, error) {
	_, err := c.storage.GetCommentByID(ctx, commentID)
	if err != nil {
		if !errors.Is(err, comment_storage.ErrCommentNotFound) {
			return nil, fmt.Errorf("GetCommentByID: %w", err)
		}
		return nil, comment_storage.ErrCommentNotFound
	}

	replies, err := c.storage.GetReplies(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("GetReplies: %w", err)
	}

	return replies, nil
}

func (c *CommentService) CreateComment(ctx context.Context, comment model.CreateCommentInput) (string, error) {
	if len(comment.Text) > maxTextLength {
		return "", ErrTextIsTooLarge
	}
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
	pst, err := c.postStorage.GetPostByID(ctx, postID)
	if err != nil {
		if !errors.Is(err, post_storage.ErrPostNotFound) {
			return "", fmt.Errorf("post_storage.GetPostByID: %w", err)
		}
		return "", post_storage.ErrPostNotFound
	}
	if !pst.CommentsAllowed {
		return "", fmt.Errorf("comments are not allowed in post: %s", pst.ID)
	}
	var comm = model.Comment{
		ID:       newID.String(),
		PostID:   comment.PostID,
		Text:     comment.Text,
		AuthorID: comment.AuthorID,
	}

	_, err = c.storage.GetCommentByID(ctx, comment.ReplyTo)
	if err != nil {
		if !errors.Is(err, comment_storage.ErrCommentNotFound) {
			return "", fmt.Errorf("comment_storage.GetCommentByID: %w", err)
		}
		return "", comment_storage.ErrCommentNotFound
	}

	comm.ReplyTo = &comment.ReplyTo
	id, err := c.storage.CreateComment(ctx, comm)
	if err != nil {
		return "", fmt.Errorf("comment_storage.CreateComment: %w", err)
	}
	var notificationPayload model.NotificationPayload
	notificationPayload.CommentAuthorID = comment.AuthorID
	notificationPayload.Text = comment.Text
	notificationPayload.PostID = postID
	c.postService.NotifyAll(ctx, notificationPayload)

	return id, nil
}

func (c *CommentService) GetComments(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error) {
	if _, err := c.postStorage.GetPostByID(ctx, *filter.PostID); err != nil {
		if !errors.Is(err, post_storage.ErrPostNotFound) {
			return nil, fmt.Errorf("post_storage.GetPostByID: %w", err)
		}
		return nil, post_storage.ErrPostNotFound
	}

	comms, err := c.storage.GetCommentsByPostID(ctx, filter)
	if err != nil {
		return nil, err
	}
	if comms == nil {
		comms = []model.Comment{}
	}
	return comms, nil
}

func (c *CommentService) GetCommentsByAuthorID(ctx context.Context, authorID string) ([]model.Comment, error) {
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

func NewCommentService(storage commentStorage, authorStorage authorStorage, postStorage postStorage, postService postService) *CommentService {
	return &CommentService{storage: storage, authorStorage: authorStorage, postStorage: postStorage, postService: postService}
}
