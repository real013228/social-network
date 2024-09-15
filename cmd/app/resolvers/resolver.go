package resolvers

import (
	"context"
	"github.com/real013228/social-network/internal/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userService userService
}

type userService interface {
	CreateUser(ctx context.Context, user model.CreateUserInput) error
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, filter model.UsersFilter) (model.User, error)
}
