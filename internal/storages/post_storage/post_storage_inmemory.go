package post_storage

import (
	"context"
	"errors"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/tools"
	"strconv"
	"sync"
)

var (
	ErrInvalidPaginationParams = errors.New("invalid pagination params")
)

type PostStorageInMemory struct {
	posts []model.Post
	cnt   int
	mu    sync.RWMutex
}

func NewPostStorageInMemory() *PostStorageInMemory {
	return &PostStorageInMemory{
		posts: make([]model.Post, 0),
		cnt:   0,
		mu:    sync.RWMutex{},
	}
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
