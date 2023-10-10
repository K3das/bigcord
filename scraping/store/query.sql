-- name: GetState :one
SELECT *
FROM states
WHERE channel_id = ?
LIMIT 1;

-- name: ListStates :many
SELECT *
FROM states
ORDER BY channel_id;

-- name: SetState :exec
INSERT OR
REPLACE
INTO states (guild_id, channel_id, type, name, state, message_offset)
VALUES (?, ?, ?, ?, ?, ?);

-- name: SetStateIfNotCompleted :exec
UPDATE OR IGNORE states
SET state = ?
WHERE channel_id = ?
  AND state != 4;

-- name: ClearState :exec
DELETE
FROM states
WHERE channel_id = ?;

-- name: ClearCompletedStates :exec
DELETE
FROM states
WHERE state = 4;

-- name: ClearAll :exec
-- noinspection SqlWithoutWhere
DELETE FROM states;
VACUUM;
