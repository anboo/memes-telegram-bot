-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN age SMALLINT DEFAULT NULL;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN sex CHAR(1) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN age;
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN sex;
-- +goose StatementEnd
