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
	PostID     *string `json:"authorID,omitempty"`
	PageLimit  int     `json:"pageLimit"`
	PageNumber int     `json:"pageNumber"`
}

type CreateCommentInput struct {
	Text     string `json:"text"`
	PostID   string `json:"postID"`
	AuthorID string `json:"authorID"`
}

type CreateCommentPayload struct {
	CommentID string `json:"comment"`
}
