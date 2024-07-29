-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_id BIGINT PRIMARY KEY,
    client_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    stored_until TIMESTAMP NOT NULL,
    taken_at TIMESTAMP,
    takeback_time TIMESTAMP,
    returned BOOLEAN DEFAULT FALSE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
