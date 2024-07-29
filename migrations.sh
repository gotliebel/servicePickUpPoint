goose -dir ./db/migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" status

goose -dir ./db/migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up