package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/winterochek/todo-server/internal/app/store"
)

type Store struct {
	db *sql.DB
	ur *UserRepository
	tr *TaskRepository
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

func (s *Store) Task() store.TaskRepository {
	if s.tr != nil {
		return s.tr
	}
	s.tr = &TaskRepository{
		store: s,
	}

	return s.tr
}
