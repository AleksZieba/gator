-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
  )
  RETURNING *
)

SELECT inserted_feed_follows.*, f.name AS feed_name, u.name AS user_name
From inserted_feed_follows
INNER JOIN feeds f ON inserted_feed_follows.feed_id = f.id
INNER JOIN users u ON inserted_feed_follows.user_id = u.id;

-- name: GetFeedFollowsForUser :many
SELECT f.name AS feed_name, u.name AS user_name
FROM feed_follows ff
INNER JOIN feeds f ON ff.feed_id = f.id
INNER JOIN users u ON ff.user_id = u.id
WHERE ff.user_id = $1;

-- name: DeleteFollow :exec
DELETE From feed_follows ff
USING feeds f, users u
WHERE ff.user_id = u.id AND ff.feed_id = f.id
AND ff.user_id = $1
AND f.url = $2;