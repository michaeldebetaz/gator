-- name: CreateFeed :one
INSERT INTO
    feeds (name, url, user_id)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: GetAllFeeds :many
SELECT
    sqlc.embed(feeds),
    sqlc.embed(users)
FROM
    feeds
    INNER JOIN users ON feeds.user_id = users.id;
