-- +goose Up
-- +goose StatementBegin
CREATE TABLE reserved_memes (
    mem_id UUID,
    user_id UUID
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reserved_memes;
-- +goose StatementEnd
