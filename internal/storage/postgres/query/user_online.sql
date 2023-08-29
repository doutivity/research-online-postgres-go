-- name: UserOnlineUpsert :exec
INSERT INTO user_online (user_id, online)
VALUES (@user_id, @online)
ON CONFLICT (user_id) DO UPDATE
    SET online = @online;

-- name: UserOnlineUpdate :exec
UPDATE user_online
SET online = @online
WHERE user_id = @user_id;

-- name: UserOnlineBatchExecUpsert :batchexec
INSERT INTO user_online (user_id, online)
VALUES (@user_id, @online)
ON CONFLICT (user_id) DO UPDATE
    SET online = @online;

-- name: UserOnlineBatchExecUpdate :batchexec
UPDATE user_online
SET online = @online
WHERE user_id = @user_id;

-- name: UserOnlineUnnestUpsert :exec
INSERT INTO user_online (user_id, online)
VALUES (unnest(@user_ids::BIGINT[]),
        unnest(@onlines::TIMESTAMP[]))
ON CONFLICT (user_id) DO UPDATE
    SET online = excluded.online;

-- name: UserOnlineUnnestUpdate :exec
UPDATE user_online AS to_t
SET online = from_t.online
FROM (
         SELECT unnest(@user_ids::BIGINT[])   AS user_id,
                unnest(@onlines::TIMESTAMP[]) AS online
     ) AS from_t
WHERE to_t.user_id = from_t.user_id;
