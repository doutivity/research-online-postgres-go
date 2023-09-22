-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_online
(
    user_id BIGINT PRIMARY KEY,
    online  TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_online;
-- +goose StatementEnd
