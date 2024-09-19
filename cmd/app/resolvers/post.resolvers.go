package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.50

import (
	"context"

	"github.com/real013228/social-network/graph"
	"github.com/real013228/social-network/internal/model"
	"github.com/real013228/social-network/tools"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.CreatePostPayload, error) {
	var createPostPayload model.CreatePostPayload
	postId, err := r.postService.CreatePost(ctx, input)
	if err != nil {
		return nil, err
	}
	posts, err := r.postService.GetPostsByFilter(ctx, model.PostsFilter{PostID: &postId})
	if err != nil || len(posts) < 1 {
		return nil, err
	}
	createPostPayload.Post = &posts[0]
	return &createPostPayload, nil
}

// Subscribe is the resolver for the subscribe field.
func (r *mutationResolver) Subscribe(ctx context.Context, input model.SubscribeInput) (*model.SubscribePayload, error) {
	payload, err := r.postService.Subscribe(ctx, input)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *model.Post) ([]*model.Comment, error) {
	comms, err := r.commentService.GetComments(ctx, model.CommentsFilter{PostID: &obj.ID, PageLimit: tools.DefaultPageLimit, PageNumber: tools.DefaultPageNumber})
	if err != nil {
		return nil, err
	}

	var res []*model.Comment
	for _, comm := range comms {
		comm := comm
		res = append(res, &comm)
	}
	return res, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, filter *model.PostsFilter) (*model.PostPayload, error) {
	posts, err := r.postService.GetPosts(ctx, *filter)
	if err != nil {
		return nil, err
	}

	var res []*model.Post
	for _, post := range posts {
		post := post
		res = append(res, &post)
	}
	var payload model.PostPayload
	payload.Posts = res
	return &payload, nil
}

// Post returns graph.PostResolver implementation.
func (r *Resolver) Post() graph.PostResolver { return &postResolver{r} }

type postResolver struct{ *Resolver }
