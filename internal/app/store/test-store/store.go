package teststore

import (
	"github.com/winterochek/todo-server/internal/app/model"
	"github.com/winterochek/todo-server/internal/app/store"
)

type Store struct {
	ur *UserRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.ur != nil {
		return s.ur
	}

	s.ur = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.ur
}