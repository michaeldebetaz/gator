-- name: CreatePost :one
INSERT INTO
    posts (title, url, description, published_at, feed_id)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: GetPostByUrl :one
SELECT
    *
FROM
    posts
WHERE
    url = $1;

-- name: GetPostsForUser :many
SELECT
    sqlc.embed(posts),
    sqlc.embed(users)
FROM
    posts
    INNER JOIN feeds ON posts.feed_id = feeds.id
    INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id
    INNER JOIN users ON feed_follows.user_id = users.id
WHERE
    users.id = $1
ORDER BY
    posts.published_at DESC
LIMIT
    $2;
