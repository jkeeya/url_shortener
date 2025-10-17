-- +goose Up
-- +goose StatementBegin
CREATE TABLE short_links (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(7)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE short_links;
-- +goose StatementEnd
