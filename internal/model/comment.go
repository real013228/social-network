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
	CommentID *string `json:"commentID,omitempty"`
	PostID    *string `json:"postID,omitempty"`
	AuthorID  *string `json:"authorID,omitempty"`
}

type CreateCommentInput struct {
	Text     string `json:"text"`
	PostID   string `json:"postID"`
	AuthorID string `json:"authorID"`
}

type CreateCommentPayload struct {
	Comment *Comment `json:"comment"`
}
