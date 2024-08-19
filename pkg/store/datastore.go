package store

import (
	"errors"

	"github.com/ssych/file_service/pkg/storage"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	storage *storage.DB
}

func NewStore(db *storage.DB) *Store {
	return &Store{storage: db}
}
