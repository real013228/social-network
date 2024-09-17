package post_service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages/user_storage"
)

const (
	DefaultPageLimit  = 10
	DefaultPageNumber = 0
)

var (
	ErrAuthorNotExist = errors.New("author doesn't exist")
)

type postStorage interface {
	CreatePost(ctx context.Context, post model.Post) (string, error)
	GetPosts(ctx context.Context, filter model.PostsFilter) ([]model.Post, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error)
	GetPostByID(ctx context.Context, postID string) (model.Post, error)
}

type userStorage interface {
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
}

type PostService struct {
	storage     postStorage
	userStorage userStorage
}

func NewPostService(storage postStorage, userStorage userStorage) *PostService {
	return &PostService{storage: storage, userStorage: userStorage}
}

func (s *PostService) CreatePost(ctx context.Context, post model.CreatePostInput) (string, error) {
	newID := uuid.New()
	authorID := post.AuthorID
	_, err := s.userStorage.GetUserByID(ctx, model.UsersFilter{UserID: &authorID})
	if err != nil {
		if !errors.Is(err, user_storage.ErrUserNotFound) {
			return "", err
		}
		return "", ErrAuthorNotExist
	}
	var pst model.Post = model.Post{
		ID:              newID.String(),
		Title:           post.Title,
		Description:     post.Description,
		AuthorID:        post.AuthorID,
		CommentsAllowed: post.CommentsAllowed,
	}
	id, err := s.storage.CreatePost(ctx, pst)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *PostService) GetPosts(ctx context.Context, filter model.PostsFilter) ([]model.Post, error) {
	posts, err := s.GetPosts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) GetPostsByFilter(ctx context.Context, filter model.PostsFilter) ([]model.Post, error) {
	var posts []model.Post
	if filter.PostID != nil {
		post, err := s.storage.GetPostByID(ctx, *filter.PostID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
		return posts, nil
	}
	if filter.AuthorID != nil {
		postsByUserID, err := s.storage.GetPostsByUserID(ctx, *filter.AuthorID)
		if err != nil {
			return nil, err
		}
		if filter.WithComments != nil {
			if *filter.WithComments {
				// todo join comments to posts
			} else {
				return postsByUserID, nil
			}
		}
		return posts, nil
	}

	return posts, nil
}
