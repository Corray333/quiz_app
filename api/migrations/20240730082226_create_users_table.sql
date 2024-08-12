-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    password TEXT NOT NULL DEFAULT '',
    role VARCHAR(16) NOT NULL DEFAULT 'user',
    tg_id BIGINT NOT NULL DEFAULT 0,
    username VARCHAR(64) NOT NULL DEFAULT '',
    email VARCHAR(64) NOT NULL DEFAULT '',
    phone VARCHAR(64) NOT NULL DEFAULT '',
    data JSONB NOT NULL DEFAULT '{}'::JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
