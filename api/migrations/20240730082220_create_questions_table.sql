-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS questions (
    quiz_id BIGINT NOT NULL REFERENCES quizzes (quiz_id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    type VARCHAR(32) NOT NULL,
    next_question_id INTEGER,
    data JSONB NOT NULL,
    FOREIGN KEY (next_question_id) REFERENCES questions (question_id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS questions;
-- +goose StatementEnd
