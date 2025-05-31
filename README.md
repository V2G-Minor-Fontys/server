# Secret files

Obtain `.env` and `development.json` and place them in `configs/`

# Running

`docker compose up` - database
`go run cmd/app/main.go` - server

# Database

Use [migrate](https://github.com/golang-migrate/migrate) to set up the database:
`migrate -database postgres://testUser:<user>!@localhost:<port>/<db>?sslmode=disable -path database/migrations up`

Create a `.sql` file in `database/queries` and then run `sqlc generate` to create an SQL schema in `internal/repository`
