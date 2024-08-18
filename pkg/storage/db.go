package storage

import (
	"context"
	"database/sql"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // get posgres driver file
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/ssych/file_service/pkg/config"
)

type DB struct {
	DB *pgxpool.Pool
}

func NewDB(
	ctx context.Context,
	option *config.DBOption,
) (*DB, error) {
	db := &DB{}

	poolConfig, err := pgxpool.ParseConfig(option.ConnectString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DATABASE_URL: %v", err)
	}

	db.DB, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	return db, nil
}

func Migrate(
	option *config.DBOption,
) error {
	db, err := sql.Open("postgres", option.ConnectString)
	if err != nil {
		return err
	}

	driver, err := migrate_postgres.WithInstance(db, &migrate_postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
