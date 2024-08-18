package store

import (
	"github.com/ssych/file_service/pkg/storage"
)

type Store struct {
	storage *storage.DB
}

func NewStore(db *storage.DB) *Store {
	return &Store{storage: db}
}
