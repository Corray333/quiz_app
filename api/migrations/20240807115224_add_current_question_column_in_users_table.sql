-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN current_question BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN current_question;
-- +goose StatementEnd
