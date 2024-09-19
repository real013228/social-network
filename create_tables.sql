DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       username VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL
);

CREATE TABLE posts (
                       id UUID PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       author_id UUID REFERENCES users(id) ON DELETE SET NULL,
                       comments_allowed BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE comments (
                          id UUID PRIMARY KEY,
                          text TEXT NOT NULL,
                          post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
                          author_id UUID REFERENCES users(id) ON DELETE SET NULL,
                          reply_to UUID REFERENCES comments(id) ON DELETE SET NULL
);

CREATE TABLE subscriptions (
                           user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                           post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
                           PRIMARY KEY (user_id, post_id)
);

CREATE TABLE notifications (
                               id UUID PRIMARY KEY,
                               receiver_id UUID REFERENCES users(id) ON DELETE CASCADE,
                               text VARCHAR(255) NOT NULL,
                               post_id UUID NOT NULL REFERENCES posts(id),
                               comment_author_id UUID NOT NULL REFERENCES users(id)
);