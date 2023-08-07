-- name: UserOnlineNew :exec
INSERT INTO user_online (user_id, online)
VALUES (@user_online, @online)
ON CONFLICT (user_id) DO UPDATE
    SET online = @online;
