-- +goose Up
-- +goose StatementBegin
CREATE TABLE reserved_mem_users (
    mem_id UUID,
    user_id UUID
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX uniq_reserved_mem_users_mem_id_user_id  ON reserved_mem_users (mem_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reserved_mem_users;
-- +goose StatementEnd
