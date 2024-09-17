package model

type CreatePostInput struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	AuthorID        string `json:"authorID"`
	CommentsAllowed bool   `json:"commentsAllowed"`
}

type CreatePostPayload struct {
	Post *Post `json:"post"`
}

type Post struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Comments        []*Comment `json:"comments"`
	AuthorID        string     `json:"authorID"`
	CommentsAllowed bool       `json:"commentsAllowed"`
}

type PostPayload struct {
	Posts []*Post `json:"posts,omitempty"`
}

type PostsFilter struct {
	PostID       *string `json:"postID,omitempty"`
	AuthorID     *string `json:"authorID,omitempty"`
	WithComments *bool   `json:"withComments,omitempty"`
	PageLimit    int     `json:"pageLimit"`
	PageNumber   int     `json:"pageNumber"`
}
