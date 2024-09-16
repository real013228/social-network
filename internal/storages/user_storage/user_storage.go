package user_storage

import (
	"context"
	"errors"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/internal/storages"
)

var (
	ErrUserNotFound = errors.New("user is not found")
)

type UserStoragePostgres struct {
	client storages.Client
}

func (s UserStoragePostgres) CreateUser(ctx context.Context, user model.User) (string, error) {
	q := `
		INSERT INTO users (id, username, email) 
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var userId string
	if err := s.client.QueryRow(ctx, q, user.ID, user.Username, user.Email).Scan(&userId); err != nil {
		return "", err
	}
	return userId, nil
}

func (s UserStoragePostgres) GetUsers(ctx context.Context) ([]model.User, error) {
	q := `
		SELECT id, username, email FROM public.users
	`
	rows, err := s.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s UserStoragePostgres) GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error) {
	q := `
		SELECT id, username, email FROM public.users
		WHERE id = $1
	`
	var user model.User
	if err := s.client.QueryRow(ctx, q, filter.UserID).Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s UserStoragePostgres) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	q := `
		SELECT id, username, email FROM public.users
		WHERE email = $1
	`
	var user model.User
	if err := s.client.QueryRow(ctx, q, email).Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return model.User{}, ErrUserNotFound
	}

	return user, nil
}

func NewUserStoragePostgres(client storages.Client) *UserStoragePostgres {
	return &UserStoragePostgres{client: client}
}
