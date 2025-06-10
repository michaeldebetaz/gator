-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON UPDATE CASCADE ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE posts;
