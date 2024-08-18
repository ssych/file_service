package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateAsset(ctx context.Context, name string, uid int64, data []byte) error {
	query := `INSERT INTO assets (name, uid, data) VALUES (@name, @uid, @data)`

	args := pgx.NamedArgs{
		"name": name,
		"uid":  uid,
		"data": data,
	}

	_, err := s.storage.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}

	return nil
}

func (s *Store) AssetFindByName(ctx context.Context, name string, uid int64) ([]byte, error) {
	var data []byte
	if err := s.storage.DB.QueryRow(ctx, "select data from assets where name = $1 and uid = $2 limit 1", name, uid).Scan(&data); err != nil {
		return nil, err
	}

	return data, nil
}
