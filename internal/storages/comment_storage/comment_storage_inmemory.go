package comment_storage

import (
	"context"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/tools"
	"strconv"
	"sync"
)

type CommentStorageInMemory struct {
	comments []model.Comment
	cnt      int
	mu       sync.RWMutex
}

func NewCommentStorageInMemory() *CommentStorageInMemory {
	return &CommentStorageInMemory{
		comments: make([]model.Comment, 0),
		cnt:      0,
		mu:       sync.RWMutex{},
	}
}

func (c *CommentStorageInMemory) CreateComment(ctx context.Context, input model.Comment) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	input.ID = strconv.Itoa(c.cnt)
	c.cnt = c.cnt + 1
	c.comments = append(c.comments, input)
	return input.ID, nil
}

func (c *CommentStorageInMemory) GetCommentsByPostID(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]model.Comment, 0)
	for _, comment := range c.comments {
		comment := comment
		if comment.PostID == *filter.PostID {
			res = append(res, comment)
		}
	}

	startIndex, endIndex, err := tools.Paginate(filter.PageLimit, filter.PageNumber, c.cnt)
	if err != nil {
		return nil, err
	}
	return res[startIndex:endIndex], nil
}

func (c *CommentStorageInMemory) GetCommentsByUserID(ctx context.Context, userID string) ([]model.Comment, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]model.Comment, 0)
	for _, comment := range c.comments {
		comment := comment
		if comment.AuthorID == userID {
			res = append(res, comment)
		}
	}
	return res, nil
}
