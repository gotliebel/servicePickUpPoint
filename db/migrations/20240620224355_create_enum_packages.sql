-- +goose Up
-- +goose StatementBegin
CREATE TYPE packages as enum ('bag', 'box', 'wrapping');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS packages;
-- +goose StatementEnd
