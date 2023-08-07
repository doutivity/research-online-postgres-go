-- name: UserOnlineNew :exec
INSERT INTO user_online (user_id, online)
VALUES (@user_id, @online)
ON CONFLICT (user_id) DO UPDATE
    SET online = @online;

-- name: UserOnlineAll :many
SELECT user_id, online
FROM user_online
ORDER BY user_id;
