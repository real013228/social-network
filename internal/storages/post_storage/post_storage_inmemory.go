package post_storage

import (
	"context"
	"errors"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/tools"
	"sync"
)

var (
	ErrInvalidPaginationParams = errors.New("invalid pagination params")
)

type userStorage interface {
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
}

type PostStorageInMemory struct {
	posts       map[string]model.Post
	mu          sync.RWMutex
	subscribers map[string][]string
	userStorage userStorage
}

func (p *PostStorageInMemory) Subscribe(ctx context.Context, subscribeInput model.SubscribeInput) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := p.subscribers[subscribeInput.PostID]
	if ok {
		p.subscribers[subscribeInput.PostID] = append(p.subscribers[subscribeInput.PostID], subscribeInput.UserID)
	} else {
		p.subscribers[subscribeInput.PostID] = []string{subscribeInput.UserID}
	}
	return "successfully subscribed to post", nil
}

func (p *PostStorageInMemory) GetSubscribers(ctx context.Context, postID string) ([]model.User, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var res []model.User
	for _, usr := range p.subscribers[postID] {
		usrByID, err := p.userStorage.GetUserByID(ctx, model.UsersFilter{
			UserID:     &usr,
			PageLimit:  tools.DefaultPageLimit,
			PageNumber: tools.DefaultPageNumber,
		})
		if err != nil {
			return nil, err
		}
		res = append(res, usrByID)
	}
	return res, nil
}

func (p *PostStorageInMemory) CreatePost(ctx context.Context, post model.Post) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.posts[post.ID] = post
	return post.ID, nil
}

func (p *PostStorageInMemory) GetPosts(ctx context.Context, filter model.PostsFilter) ([]model.Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	startIndex, endIndex, err := tools.Paginate(filter.PageLimit, filter.PageNumber, len(p.posts))
	if err != nil {
		return nil, err
	}
	cnt := 0
	var results []model.Post
	for k, _ := range p.posts {
		if cnt >= endIndex-startIndex+1 {
			break
		}
		results = append(results, p.posts[k])
		cnt++
	}
	return results, nil
}

func (p *PostStorageInMemory) GetPostsByUserID(ctx context.Context, userID string) ([]model.Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var posts []model.Post
	for _, post := range p.posts {
		if post.AuthorID == userID {
			posts = append(posts, post)
			break
		}
	}
	return posts, nil
}

func (p *PostStorageInMemory) GetPostByID(ctx context.Context, postID string) (model.Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	post, ok := p.posts[postID]
	if !ok {
		return model.Post{}, ErrPostNotFound
	}
	return post, nil
}

func NewPostStorageInMemory(userStorage userStorage) *PostStorageInMemory {
	return &PostStorageInMemory{
		posts:       make(map[string]model.Post, 0),
		mu:          sync.RWMutex{},
		subscribers: make(map[string][]string),
		userStorage: userStorage,
	}
}
