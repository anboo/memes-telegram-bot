-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN full_name VARCHAR(255) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN full_name;
-- +goose StatementEnd