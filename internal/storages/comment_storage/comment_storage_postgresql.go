package comment_storage

import (
	"context"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages"
)

type CommentStoragePostgres struct {
	client storages.Client
}

func (s *CommentStoragePostgres) CreateComment(ctx context.Context, input model.Comment) (string, error) {
	q := `
		INSERT INTO comments (id, text, post_id, author_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var id string
	if err := s.client.QueryRow(ctx, q, input.ID, input.Text, input.PostID, input.AuthorID).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (s *CommentStoragePostgres) GetCommentsByPostID(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error) {
	q := `
		SELECT id, text, post_id, author_id
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
		err = rows.Scan(&comment.ID, &comment.Text, &comment.PostID, &comment.AuthorID)
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
		SELECT id, text, post_id, author_id
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
		err = rows.Scan(&comment.ID, &comment.Text, &comment.PostID, &comment.AuthorID)
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
