-- +goose Up
-- +goose StatementBegin
CREATE TABLE memes (
   id UUID NOT NULL,
   external_id VARCHAR(255) NOT NULL,
   text TEXT,
   source VARCHAR(255) NOT NULL,
   source_from VARCHAR(255) NOT NULL,
   img VARCHAR(255) NOT NULL,
   rating INT DEFAULT 0,
   PRIMARY KEY(id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE UNIQUE INDEX uniq_memes_source_external_id ON memes (source, external_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uniq_memes_source_external_id;
DROP TABLE memes;
-- +goose StatementEnd
