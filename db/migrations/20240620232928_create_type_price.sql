-- +goose Up
-- +goose StatementBegin
CREATE TYPE price AS (
    number NUMERIC,
    currency_code TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS price;
-- +goose StatementEnd
