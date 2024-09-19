package post_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages/post_storage"
	"github.com/real013228/social-network/internal/storages/user_storage"
	"strings"
	"sync"
)

var (
	ErrAuthorNotExist = errors.New("author doesn't exist")
)

type postStorage interface {
	CreatePost(ctx context.Context, post model.Post) (string, error)
	GetPosts(ctx context.Context, filter model.PostsFilter) ([]model.Post, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error)
	GetPostByID(ctx context.Context, postID string) (model.Post, error)
	Subscribe(ctx context.Context, subscribeInput model.SubscribeInput) (string, error)
	GetSubscribers(ctx context.Context, postID string) ([]model.User, error)
}

type commentStorage interface {
	GetCommentsByPostID(ctx context.Context, filter model.CommentsFilter) ([]model.Comment, error)
}

type userStorage interface {
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
}

type userService interface {
	Notify(ctx context.Context, userID string, payload model.NotificationPayload)
}

type PostService struct {
	storage        postStorage
	userStorage    userStorage
	userService    userService
	commentStorage commentStorage
}

func NewPostService(storage postStorage, userStorage userStorage, userService userService, commentStorage commentStorage) *PostService {
	return &PostService{storage: storage, userStorage: userStorage, userService: userService, commentStorage: commentStorage}
}

func (s *PostService) Subscribe(ctx context.Context, subscribeInput model.SubscribeInput) (model.SubscribePayload, error) {
	msg, err := s.storage.Subscribe(ctx, subscribeInput)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return model.SubscribePayload{}, fmt.Errorf("user already subscribed to this post, userID: %s, postID: %s", subscribeInput.UserID, subscribeInput.PostID)
		}
		return model.SubscribePayload{}, err
	}
	subs, err := s.storage.GetSubscribers(ctx, subscribeInput.PostID)
	if err != nil {
		return model.SubscribePayload{}, err
	}
	for _, sub := range subs {
		if sub.ID == subscribeInput.UserID {
			return model.SubscribePayload{}, fmt.Errorf("user already subscribed to this post, userID: %s, postID: %s", subscribeInput.UserID, subscribeInput.PostID)
		}
	}
	var payload model.SubscribePayload
	payload.Message = &msg
	payload.Success = true
	return payload, nil
}

func (s *PostService) GetPostByID(ctx context.Context, postID string) (model.Post, error) {
	post, err := s.storage.GetPostByID(ctx, postID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return model.Post{}, post_storage.ErrPostNotFound
		}
		return model.Post{}, err
	}
	return post, nil
}
func (s *PostService) CreatePost(ctx context.Context, post model.CreatePostInput) (string, error) {
	newID := uuid.New()
	authorID := post.AuthorID
	_, err := s.userStorage.GetUserByID(ctx, model.UsersFilter{UserID: &authorID})
	if err != nil {
		if !errors.Is(err, user_storage.ErrUserNotFound) && !strings.Contains(err.Error(), "no rows in result set") {
			return "", err
		}
		return "", ErrAuthorNotExist
	}
	var pst = model.Post{
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
	posts, err := s.storage.GetPosts(ctx, filter)
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
		if filter.WithComments != nil && *filter.WithComments {
			var commentsFilter model.CommentsFilter
			commentsFilter.PostID = &post.ID
			commentsFilter.PageLimit = filter.PageLimit
			commentsFilter.PageNumber = filter.PageNumber
			comms, err := s.commentStorage.GetCommentsByPostID(ctx, commentsFilter)
			if err != nil {
				return nil, fmt.Errorf("GetCommentsByPostID: %w", err)
			}
			var res []*model.Comment
			for _, comment := range comms {
				comment := comment
				res = append(res, &comment)
			}
			post.Comments = res
		}
		posts = append(posts, post)
		return posts, nil
	}
	if filter.AuthorID != nil {
		postsByUserID, err := s.storage.GetPostsByUserID(ctx, *filter.AuthorID)
		if err != nil {
			return nil, err
		}
		if filter.WithComments != nil && *filter.WithComments {
			for _, post := range postsByUserID {
				comms, err := s.commentStorage.GetCommentsByPostID(ctx, model.CommentsFilter{PostID: &post.ID})
				if err != nil {
					return nil, fmt.Errorf("GetCommentsByPostID: %w", err)
				}
				var res []*model.Comment
				for _, comment := range comms {
					comment := comment
					res = append(res, &comment)
				}
				post.Comments = res
			}
		}
		return posts, nil
	}

	return posts, nil
}

func (s *PostService) NotifyAll(ctx context.Context, payload model.NotificationPayload) {
	subscribers, err := s.storage.GetSubscribers(ctx, payload.PostID)
	if err == nil {
		var wg sync.WaitGroup
		wg.Add(len(subscribers))
		for _, sub := range subscribers {
			go func(ctx context.Context, subID string, payload model.NotificationPayload) {
				s.userService.Notify(ctx, subID, payload)
				wg.Done()
			}(ctx, sub.ID, payload)
		}
		wg.Wait()
	}
}
