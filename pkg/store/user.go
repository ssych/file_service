package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID           int64
	Login        string
	PasswordHash string

	CreatedAt time.Time
}

func (s *Store) UserFindByLogin(ctx context.Context, login string) (*User, error) {
	rows, err := s.storage.DB.Query(ctx, "select * from users where login = $1 limit 1", login)
	if err != nil {
		return nil, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
	}

	return &users[0], nil
}
