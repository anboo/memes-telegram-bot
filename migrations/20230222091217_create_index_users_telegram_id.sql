-- +goose Up
-- +goose StatementBegin
CREATE INDEX users_telegram_id ON users (telegram_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX users_telegram_id;
-- +goose StatementEnd
