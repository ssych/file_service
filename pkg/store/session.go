package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Session struct {
	ID        string
	UID       int64
	CreatedAt time.Time
}

func (s *Store) CreateSession(ctx context.Context, userID int64) (string, error) {
	query := `INSERT INTO sessions (id, uid) VALUES (@id, @uid)`
	token := generateToken(20)

	args := pgx.NamedArgs{
		"id":  token,
		"uid": userID,
	}

	_, err := s.storage.DB.Exec(ctx, query, args)
	if err != nil {
		return "", fmt.Errorf("unable to insert row: %w", err)
	}

	return token, nil
}

func (s *Store) SessionFindByID(ctx context.Context, id string) (*Session, error) {
	rows, err := s.storage.DB.Query(ctx, "select * from sessions where id = $1 limit 1", id)
	if err != nil {
		return nil, err
	}

	sessions, err := pgx.CollectRows(rows, pgx.RowToStructByName[Session])
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return nil, ErrNotFound
	}

	return &sessions[0], nil
}

func generateToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
