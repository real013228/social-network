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
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Posts    []*Post `json:"posts"`
}

type UserPayload struct {
	Users []*User `json:"users,omitempty"`
}

type UsersFilter struct {
	UserID *string `json:"userID,omitempty"`
}
