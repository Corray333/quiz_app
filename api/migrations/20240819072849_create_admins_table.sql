-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS admins(
    tg_id BIGINT NOT NULL DEFAULT 0,
    username TEXT NOT NULL PRIMARY KEY
);
INSERT INTO admins(tg_id, username) VALUES(56218566, 'incetro');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS admins;
-- +goose StatementEnd
