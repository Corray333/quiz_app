-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS quizzes (
    quiz_id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(256) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    cover TEXT NOT NULL DEFAULT '',
    type VARCHAR(16) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS quizzes;
-- +goose StatementEnd
