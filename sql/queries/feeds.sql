-- name: CreateFeed :one
INSERT INTO
    feeds (name, url, user_id)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: GetFeedByUrl :one
SELECT
    feeds.*
FROM
    feeds
WHERE
    url = $1;

-- name: GetNextFeedToFetch :one
SELECT
    *
FROM
    feeds
ORDER BY
    last_fetched_at ASC NULLS FIRST
LIMIT
    1;

-- name: MarkFeedAsFetched :exec
UPDATE
    feeds
SET
    last_fetched_at = NOW(),
    updated_at = NOW()
WHERE
    id = $1;

-- name: GetAllFeedsWithUsers :many
SELECT
    sqlc.embed(feeds),
    sqlc.embed(users)
FROM
    feeds
    INNER JOIN users ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO
        feed_follows (feed_id, user_id)
    VALUES
        ($1, $2)
    RETURNING
        *
)
SELECT
    inserted_feed_follow.*,
    sqlc.embed(feeds),
    sqlc.embed(users)
FROM
    inserted_feed_follow
    INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
    INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsByUser :many
SELECT
    sqlc.embed(feeds),
    sqlc.embed(users)
FROM
    feed_follows
    INNER JOIN feeds ON feed_follows.feed_id = feeds.id
    INNER JOIN users ON feed_follows.user_id = users.id
WHERE
    feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM
    feed_follows
WHERE
    feed_id = $1
    AND user_id = $2;
