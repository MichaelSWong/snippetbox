# Frontend Masters: Go For Developers

## The path to interact with the Database

1. Start with Database Layer
2. Hook the Database Layer to the API Handler
3. Then Route It

Store -> Handler -> Route

## Log into Postgres

- psql -U postgres -h localhost -p 5432

## Migrations:
goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
