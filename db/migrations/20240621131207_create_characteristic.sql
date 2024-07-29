-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS characteristic (
    order_id BIGINT PRIMARY KEY,
    package packages NOT NULL,
    weight FLOAT NOT NULL,
    price price NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS characteristic;
-- +goose StatementEnd
