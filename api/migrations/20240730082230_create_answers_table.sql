-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS answers (
    question_id INTEGER REFERENCES questions (question_id),
    user_id BIGINT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    answer JSONB NOT NULL DEFAULT '[]'::JSONB,
    checked BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (question_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS answers;
-- +goose StatementEnd
