-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID NOT NULL,
    telegram_id VARCHAR(255) NOT NULL,
    username VARCHAR(255) DEFAULT NULL,
    created_at timestamp,
    PRIMARY KEY(id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX uniq_users_telegram_id ON users (telegram_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uniq_users_telegram_id;
DROP TABLE users;
-- +goose StatementEnd
