-- +goose up
CREATE TABLE posts(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    UNIQUE(url),
    CONSTRAINT fk_feed_id
        FOREIGN KEY (feed_id)
        REFERENCES feeds(id)
        ON DELETE CASCADE
);

-- +goose down
DELETE TABLE posts;