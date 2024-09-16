package user_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages/user_storage"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type userStorage interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type UserService struct {
	storage userStorage
}

func (s UserService) CreateUser(ctx context.Context, user model.CreateUserInput) (string, error) {
	storedUser, err := s.storage.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if !errors.Is(err, user_storage.ErrUserNotFound) {
			return "", fmt.Errorf("user_storage.GetUserByEmail email=%s: %w", user.Email, err)
		}
	}
	if storedUser.Email == user.Email {
		return "", ErrUserExists
	}

	newID := uuid.New()
	var usr = model.User{
		ID:       newID.String(),
		Email:    user.Email,
		Username: user.Username,
	}
	email, err := s.storage.CreateUser(ctx, usr)
	if err != nil {
		return "", fmt.Errorf("user_storage.CreateUser email:%s %w", user.Email, err)
	}
	return email, nil
}

func (s UserService) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.storage.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("user_storage.GetUsers %w", err)
	}
	return users, nil
}

func (s UserService) GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error) {
	user, err := s.storage.GetUserByID(ctx, filter)
	if err != nil {
		return model.User{}, fmt.Errorf("user_storage.GetUserByID %w", err)
	}

	return user, nil
}

func NewUserService(storage userStorage) *UserService {
	return &UserService{storage: storage}
}
