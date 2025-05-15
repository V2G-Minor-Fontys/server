# Secret files

Obtain `.env` and `development.json` and place them in `configs/`

# Running

`docker compose up` - database
`go run cmd/app/main.go` - server

# Database

Use [migrate](https://github.com/golang-migrate/migrate) to set up the database:
`migrate -database postgres://testUser:TestPass1234!@localhost:5432/testDb?sslmode=disable -path database/migrations up`

Edit a `.sql` file in `database/queries` and then run `sqlc generate`
