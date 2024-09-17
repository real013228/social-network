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
	GetUsers(ctx context.Context, filter model.UsersFilter) ([]model.User, error)
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	GetNotifications(ctx context.Context, filter model.UsersFilter) ([]model.NotificationPayload, error)
	notify(ctx context.Context, userID string, payload model.NotificationPayload)
}

type UserService struct {
	storage userStorage
}

func (s *UserService) GetNotifications(ctx context.Context, filter model.UsersFilter) ([]model.NotificationPayload, error) {
	notifications, err := s.storage.GetNotifications(ctx, filter)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *UserService) CreateUser(ctx context.Context, user model.CreateUserInput) (string, error) {
	storedUser, err := s.storage.GetUserByEmail(ctx, user.Email)
	//todo add email validation
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
	id, err := s.storage.CreateUser(ctx, usr)
	if err != nil {
		return "", fmt.Errorf("user_storage.CreateUser email:%s %w", user.Email, err)
	}
	return id, nil
}

func (s *UserService) GetUsers(ctx context.Context, filter model.UsersFilter) ([]model.User, error) {
	users, err := s.storage.GetUsers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("user_storage.GetUsers %w", err)
	}

	return users, nil
}

func (s *UserService) GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error) {
	user, err := s.storage.GetUserByID(ctx, filter)
	//todo add validation, when user doesnt exist
	if err != nil {
		return model.User{}, fmt.Errorf("user_storage.GetUserByID %w", err)
	}

	return user, nil
}

func (s *UserService) notify(ctx context.Context, userID string, payload model.NotificationPayload) {
	s.storage.notify(ctx, userID, payload)
}

func NewUserService(storage userStorage) *UserService {
	return &UserService{storage: storage}
}
