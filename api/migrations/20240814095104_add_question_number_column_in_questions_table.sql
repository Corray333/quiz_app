-- +goose Up
-- +goose StatementBegin
ALTER TABLE questions ADD COLUMN question_number INTEGER NOT NULL DEFAULT 0;
ALTER TABLE questions DROP COLUMN next_question_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE questions ADD COLUMN next_question_id INTEGER;
ALTER TABLE questions DROP COLUMN question_number;
-- +goose StatementEnd
