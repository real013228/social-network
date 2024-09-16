package comment_service

import (
	"context"
	"github.com/real013228/social-network/internal/model"
)

type commentStorage interface {
	CreateComment(ctx context.Context, input model.Comment) (string, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]model.Comment, error)
	GetCommentsByUserID(ctx context.Context, userID string) ([]model.Comment, error)
}
