package post_storage

import (
	"context"
	"errors"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/services/post_service"
	"github.com/real013228/social-network/tools"
	"strconv"
	"sync"
)

var (
	ErrInvalidPaginationParams = errors.New("invalid pagination params")
)

type userStorage interface {
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
}

type PostStorageInMemory struct {
	posts       []model.Post
	cnt         int
	mu          sync.RWMutex
	subscribers map[int][]int
	userStorage userStorage
}

func (p *PostStorageInMemory) Subscribe(ctx context.Context, subscribeInput model.SubscribeInput) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	pId, err := strconv.Atoi(subscribeInput.PostID)
	if err != nil {
		return "", err
	}
	uId, err := strconv.Atoi(subscribeInput.UserID)
	if err != nil {
		return "", err
	}
	_, ok := p.subscribers[pId]
	if ok {
		p.subscribers[pId] = append(p.subscribers[pId], uId)
	} else {
		p.subscribers[pId] = []int{uId}
	}
	return "successfully subscribed to post", nil
}

func (p *PostStorageInMemory) getSubscribers(ctx context.Context, postID string) ([]model.User, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pId, err := strconv.Atoi(postID)
	if err != nil {
		return nil, err
	}
	var res []model.User
	for _, usr := range p.subscribers[pId] {
		usrID := strconv.Itoa(usr)
		usrByID, err := p.userStorage.GetUserByID(ctx, model.UsersFilter{
			UserID:     &usrID,
			PageLimit:  post_service.DefaultPageLimit,
			PageNumber: post_service.DefaultPageNumber,
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
	post.ID = strconv.Itoa(p.cnt)
	p.cnt = p.cnt + 1
	p.posts = append(p.posts, post)
	return post.ID, nil
}

func (p *PostStorageInMemory) GetPosts(ctx context.Context, filter model.PostsFilter) ([]model.Post, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	startIndex, endIndex, err := tools.Paginate(filter.PageLimit, filter.PageNumber, p.cnt)
	if err != nil {
		return nil, err
	}
	return p.posts[startIndex:endIndex], nil
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
	for _, post := range p.posts {
		if post.ID == postID {
			return post, nil
		}
	}
	return model.Post{}, ErrPostNotFound
}

func NewPostStorageInMemory() *PostStorageInMemory {
	return &PostStorageInMemory{
		posts: make([]model.Post, 0),
		cnt:   0,
		mu:    sync.RWMutex{},
	}
}
