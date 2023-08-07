-- name: UserOnlineUpsert :exec
INSERT INTO user_online (user_id, online)
VALUES (@user_id, @online)
ON CONFLICT (user_id) DO UPDATE
    SET online = @online;

-- name: UserOnlineUpdate :exec
UPDATE user_online
SET online = @online
WHERE user_id = @user_id;

-- name: UserOnlineBatchUpsert :exec
INSERT INTO user_online (user_id, online)
SELECT user_id, online
FROM (
         SELECT unnest(@user_ids::BIGINT[])                  AS user_id,
                unnest(@onlines::TIMESTAMP WITH TIME ZONE[]) AS online
     ) AS from_t
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;

-- name: UserOnlineBatchUpdate :exec
UPDATE user_online AS to_t
SET online = from_t.online
FROM (
         SELECT unnest(@user_ids::BIGINT[])                  AS user_id,
                unnest(@onlines::TIMESTAMP WITH TIME ZONE[]) AS online
     ) AS from_t (user_id, online)
WHERE to_t.user_id = from_t.user_id;

-- name: UserOnlineAll :many
SELECT user_id, online
FROM user_online
ORDER BY user_id;
