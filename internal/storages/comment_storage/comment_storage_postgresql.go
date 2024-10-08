package comment_storage

import (
	"context"
	"errors"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

type CommentStoragePostgres struct {
	client storages.Client
}

func (s *CommentStoragePostgres) GetCommentByID(ctx context.Context, commentID string) (model.Comment, error) {
	q := `
		SELECT id, text, post_id, author_id, reply_to
		FROM comments
		WHERE id = $1
	`
	var comment model.Comment
	if err := s.client.QueryRow(ctx, q, commentID).Scan(&comment.ID, &comment.Text, &comment.PostID, &comment.AuthorID, &comment.ReplyTo); err != nil {
		return comment, err
	}
	return comment, nil
}

func (s *CommentStoragePostgres) GetReplies(ctx context.Context, commentID string) ([]model.Comment, error) {
	q := `
		SELECT id, text, post_id, author_id, reply_to
		FROM comments
		WHERE reply_to = $1
	`
	var comments []model.Comment
	rows, err := s.client.Query(ctx, q, commentID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.Text, &comment.PostID, &comment.AuthorID, &comment.ReplyTo)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentStoragePostgres) CreateComment(ctx context.Context, input model.Comment) (string, error) {
	var q string
	var id string
	if *input.ReplyTo == "" {
		q = `
		INSERT INTO comments (id, text, post_id, author_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
		if err := s.client.QueryRow(ctx, q, input.ID, input.Text, input.PostID, input.AuthorID).Scan(&id); err != nil {
			return "", err
		}
	} else {
		q = `
		INSERT INTO comments (id, text, post_id, author_id, reply_to)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
		if err := s.client.QueryRow(ctx, q, input.ID, input.Text, input.PostID, input.AuthorID, input.ReplyTo).Scan(&id); err != nil {
			return "", err
		}
	}

	return id, nil
}

func (s *CommentStoragePostgres) GetCommentsByPostID(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error) {
	q := `
		SELECT id, text, post_id, author_id, reply_to
		FROM comments
		WHERE post_id = $1
		LIMIT $2 OFFSET $3;
	`
	rows, err := s.client.Query(ctx, q, filter.PostID, filter.PageLimit, filter.PageNumber*filter.PageLimit)
	if err != nil {
		return nil, err
	}
	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.Text, &comment.PostID, &comment.AuthorID, &comment.ReplyTo)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentStoragePostgres) GetCommentsByUserID(ctx context.Context, userID string) ([]model.Comment, error) {
	q := `
		SELECT id, text, post_id, author_id, reply_to
		FROM comments
		WHERE author_id = $1;
	`
	rows, err := s.client.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err = rows.Scan(&comment.ID, &comment.Text, &comment.PostID, &comment.AuthorID, &comment.ReplyTo)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func NewCommentStoragePostgres(client storages.Client) *CommentStoragePostgres {
	return &CommentStoragePostgres{client: client}
}
