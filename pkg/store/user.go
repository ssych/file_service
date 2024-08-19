package store

import (
	"context"
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
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return &users[0], nil
}
