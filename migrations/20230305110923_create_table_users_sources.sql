-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_sources (
    user_id UUID,
    source VARCHAR(25),
    enabled BOOL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX ux_user_id_source ON users_sources (user_id, source);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users_sources DROP INDEX ux_user_id_source;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE users_sources;
-- +goose StatementEnd
