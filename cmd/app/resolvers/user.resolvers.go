package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.50

import (
	"context"
	"fmt"

	"github.com/real013228/social-network/graph"
	"github.com/real013228/social-network/internal/model"
)

// UserID is the resolver for the userId field.
func (r *createUserPayloadResolver) UserID(ctx context.Context, obj *model.CreateUserPayload) (string, error) {
	return obj.User, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.CreateUserPayload, error) {
	var createUserPayload model.CreateUserPayload
	userId, err := r.userService.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}
	createUserPayload.User = userId
	return &createUserPayload, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, filter *model.UsersFilter) (*model.UserPayload, error) {
	var payload = &model.UserPayload{}
	if filter != nil {
		user, err := r.userService.GetUserByID(ctx, *filter)
		if err != nil {
			return nil, err
		}
		payload.Users = append(payload.Users, &user)
	} else {
		users, err := r.userService.GetUsers(ctx)
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			payload.Users = append(payload.Users, &user)
		}
	}
	return payload, nil
}

// Posts is the resolver for the posts field.
func (r *userResolver) Posts(ctx context.Context, obj *model.User) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented: Posts - posts"))
}

// CreateUserPayload returns graph.CreateUserPayloadResolver implementation.
func (r *Resolver) CreateUserPayload() graph.CreateUserPayloadResolver {
	return &createUserPayloadResolver{r}
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

// User returns graph.UserResolver implementation.
func (r *Resolver) User() graph.UserResolver { return &userResolver{r} }

type createUserPayloadResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
