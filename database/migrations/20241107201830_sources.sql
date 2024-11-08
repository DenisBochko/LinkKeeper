-- +goose Up
-- +goose StatementBegin
CREATE TABLE sources
(
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    user_url VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sources;
-- +goose StatementEnd