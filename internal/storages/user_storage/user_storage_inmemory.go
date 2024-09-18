package user_storage

import (
	"context"
	"errors"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/tools"
	"strconv"
	"sync"
)

var (
	ErrInvalidPaginationParams = errors.New("invalid pageLimit or pageNumber")
)

type UserStorageInMemory struct {
	users         []model.User
	notifications map[int][]model.NotificationPayload
	cnt           int
	mu            sync.RWMutex
}

func (u *UserStorageInMemory) GetNotifications(ctx context.Context, filter model.UsersFilter) ([]model.NotificationPayload, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	usrID, err := strconv.Atoi(*filter.UserID)
	if err != nil {
		return nil, err
	}
	return u.notifications[usrID], nil
}

func (u *UserStorageInMemory) Notify(ctx context.Context, userID string, payload model.NotificationPayload) {
	u.mu.Lock()
	defer u.mu.Unlock()
	usrID, _ := strconv.Atoi(userID)
	if notifications, ok := u.notifications[usrID]; ok {
		notifications = append(notifications, payload)
		u.notifications[usrID] = notifications
	} else {
		u.notifications[usrID] = []model.NotificationPayload{payload}
	}
}

func (u *UserStorageInMemory) CreateUser(ctx context.Context, user model.User) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.users = append(u.users, user)
	u.cnt++
	res := strconv.Itoa(u.cnt)
	return res, nil
}

func (u *UserStorageInMemory) GetUsers(ctx context.Context, filter model.UsersFilter) ([]model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	startIndex, endIndex, err := tools.Paginate(filter.PageLimit, filter.PageNumber, u.cnt)
	if err != nil {
		return nil, err
	}
	return u.users[startIndex:endIndex], nil
}

func (u *UserStorageInMemory) GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	userID, err := strconv.Atoi(*filter.UserID)
	if err != nil {
		return model.User{}, err
	}
	if userID < 0 || userID > u.cnt {
		return model.User{}, ErrUserNotFound
	}
	return u.users[userID-1], nil
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
		users:         make([]model.User, 0),
		cnt:           0,
		mu:            sync.RWMutex{},
		notifications: make(map[int][]model.NotificationPayload),
	}
}
