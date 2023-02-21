-- +goose Up
-- +goose StatementBegin
CREATE TABLE votes (
    id UUID,
    mem_id UUID,
    user_id UUID,
    vote SMALLINT,
    created_at timestamp,
    PRIMARY KEY (id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX uniq_vote_mem_id_user_id  ON votes (mem_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE votes;
-- +goose StatementEnd
