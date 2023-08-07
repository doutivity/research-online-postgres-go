-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_online
(
    user_id BIGINT PRIMARY KEY,
    online  TIMESTAMP WITH TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_online;
-- +goose StatementEnd
