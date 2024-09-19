package user_storage

import (
	"context"
	"errors"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/tools"
	"sync"
)

var (
	ErrInvalidPaginationParams = errors.New("invalid pageLimit or pageNumber")
	ErrInvalidIdValue          = errors.New("invalid id, must be integer")
)

type UserStorageInMemory struct {
	users         map[string]model.User
	notifications map[string][]model.NotificationPayload
	cnt           int
	mu            sync.RWMutex
}

func (u *UserStorageInMemory) GetNotifications(ctx context.Context, filter model.UsersFilter) ([]model.NotificationPayload, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.notifications[*filter.UserID], nil
}

func (u *UserStorageInMemory) Notify(ctx context.Context, userID string, payload model.NotificationPayload) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if notifications, ok := u.notifications[userID]; ok {
		notifications = append(notifications, payload)
		u.notifications[userID] = notifications
	} else {
		u.notifications[userID] = []model.NotificationPayload{payload}
	}
}

func (u *UserStorageInMemory) CreateUser(ctx context.Context, user model.User) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.users[user.ID] = user
	return user.ID, nil
}

func (u *UserStorageInMemory) GetUsers(ctx context.Context, filter model.UsersFilter) ([]model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	startIndex, endIndex, err := tools.Paginate(filter.PageLimit, filter.PageNumber, u.cnt)
	if err != nil {
		return nil, err
	}
	cnt := 0
	var results []model.User
	for k, _ := range u.users {
		if cnt >= endIndex-startIndex+1 {
			break
		}
		results = append(results, u.users[k])
		cnt++
	}
	return results, nil
}

func (u *UserStorageInMemory) GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	val, ok := u.users[*filter.UserID]
	if !ok {
		return model.User{}, ErrUserNotFound
	}
	return val, nil
}

func (u *UserStorageInMemory) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for _, user := range u.users {
		if user.Email == email {
			return user, nil
		}
	}
	return model.User{}, ErrUserNotFound
}

func NewUserStorageInMemory() *UserStorageInMemory {
	return &UserStorageInMemory{
		users:         make(map[string]model.User),
		cnt:           0,
		mu:            sync.RWMutex{},
		notifications: make(map[string][]model.NotificationPayload),
	}
}
