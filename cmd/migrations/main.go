package main

import (
	"database/sql"

	migrate "github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // get posgres driver file
	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost user=dev dbname=file_service_development sslmode=disable password=dev"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	driver, err := migrate_postgres.WithInstance(db, &migrate_postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
