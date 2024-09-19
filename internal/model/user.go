package model

type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type CreateUserPayload struct {
	User string `json:"user"`
}

type Mutation struct {
}

type Query struct {
}

type User struct {
	ID            string   `json:"id"`
	Username      string   `json:"username"`
	Email         string   `json:"email"`
	Posts         []*Post  `json:"posts"`
	Notifications []string `json:"notifications"`
}

type UserPayload struct {
	Users []*User `json:"users,omitempty"`
}

type UsersFilter struct {
	UserID     *string `json:"userID,omitempty"`
	PageLimit  int     `json:"pageLimit"`
	PageNumber int     `json:"pageNumber"`
}

type NotificationPayload struct {
	ID              string `json:"id"`
	Text            string `json:"text"`
	PostID          string `json:"postID"`
	CommentAuthorID string `json:"commentAuthorID"`
}
