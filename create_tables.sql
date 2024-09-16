DROP TABLE IF EXISTS user_posts;
DROP TABLE IF EXISTS post_comments;
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
                       author_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Ссылка на автора поста
                       comments_allowed BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE comments (
                          id UUID PRIMARY KEY,
                          text TEXT NOT NULL,
                          post_id UUID REFERENCES posts(id) ON DELETE CASCADE,    -- Ссылка на пост
                          author_id UUID REFERENCES users(id) ON DELETE SET NULL  -- Ссылка на автора комментария
);

-- Для связи постов с пользователями (множество постов у одного пользователя)
CREATE TABLE user_posts (
                            user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                            post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
                            PRIMARY KEY (user_id, post_id)
);

-- Для связи комментариев с постами (множество комментариев у одного поста)
CREATE TABLE post_comments (
                               post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
                               comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
                               PRIMARY KEY (post_id, comment_id)
);
