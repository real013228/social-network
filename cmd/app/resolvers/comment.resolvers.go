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

// Author is the resolver for the author field.
func (r *commentResolver) Author(ctx context.Context, obj *model.Comment) (*model.User, error) {
	authorId := obj.AuthorID
	author, err := r.userService.GetUserByID(ctx, model.UsersFilter{UserID: &authorId})
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// CommentID is the resolver for the commentID field.
func (r *createCommentPayloadResolver) CommentID(ctx context.Context, obj *model.CreateCommentPayload) (string, error) {
	panic(fmt.Errorf("not implemented: CommentID - commentID"))
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.CreateCommentPayload, error) {
	id, err := r.commentService.CreateComment(ctx, input)
	if err != nil {
		return nil, err
	}
	var createdComment model.CreateCommentPayload
	createdComment.CommentID = id
	return &createdComment, nil
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, filter *model.CommentsFilter) (*model.CommentPayload, error) {
	comms, err := r.commentService.GetCommentsByAuthorID(ctx, *filter.AuthorID)
	if err != nil {
		return nil, err
	}

	var commPayload model.CommentPayload
	var res []*model.Comment
	for _, comm := range comms {
		comm := comm
		res = append(res, &comm)
	}
	commPayload.Comments = res

	return &commPayload, nil
}

// Comment returns graph.CommentResolver implementation.
func (r *Resolver) Comment() graph.CommentResolver { return &commentResolver{r} }

// CreateCommentPayload returns graph.CreateCommentPayloadResolver implementation.
func (r *Resolver) CreateCommentPayload() graph.CreateCommentPayloadResolver {
	return &createCommentPayloadResolver{r}
}

type commentResolver struct{ *Resolver }
type createCommentPayloadResolver struct{ *Resolver }
