-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd
