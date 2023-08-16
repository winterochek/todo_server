package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/winterochek/todo-server/internal/app/store"
)

type Store struct {
	db *sql.DB
	ur *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.ur != nil {
		return s.ur
	}
	s.ur = &UserRepository{
		store: s,
	}
	return s.ur
}