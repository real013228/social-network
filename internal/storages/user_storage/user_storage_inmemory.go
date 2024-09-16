package user_storage

import (
	"context"
	"errors"
	"github.com/real013228/social-network/internal/model"
	"strconv"
	"sync"
)

type UserStorageInMemory struct {
	users []model.User
	cnt   int
	mu    sync.RWMutex
}

var (
	ErrInvalidPaginationParams = errors.New("invalid pageLimit or pageNumber")
)

func (u *UserStorageInMemory) CreateUser(ctx context.Context, user model.User) (string, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.users = append(u.users, user)
	u.cnt++
	res := strconv.Itoa(u.cnt)
	return res, nil
}

func (u *UserStorageInMemory) GetUsers(ctx context.Context, pageLimit, pageNumber int) ([]model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	if pageLimit <= 0 || pageNumber <= 0 {
		return nil, ErrInvalidPaginationParams
	}
	startIndex := (pageNumber - 1) * pageLimit
	if startIndex >= u.cnt {
		return nil, nil
	}
	endIndex := startIndex + pageLimit
	if endIndex > u.cnt {
		endIndex = u.cnt
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
