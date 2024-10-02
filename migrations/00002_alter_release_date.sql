-- +goose Up
-- +goose StatementBegin
ALTER TABLE songs
    ALTER COLUMN release_date TYPE VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE songs
    ALTER COLUMN release_date TYPE DATE;
-- +goose StatementEnd
