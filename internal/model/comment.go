package model

type Comment struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	PostID   string `json:"postID"`
	AuthorID string `json:"authorID"`
}

type CommentPayload struct {
	Comments []*Comment `json:"comments,omitempty"`
}

type CommentsFilter struct {
	AuthorID *string `json:"authorID,omitempty"`
}

type CreateCommentInput struct {
	Text     string `json:"text"`
	PostID   string `json:"postID"`
	AuthorID string `json:"authorID"`
}

type CreateCommentPayload struct {
	CommentID string `json:"comment"`
}
