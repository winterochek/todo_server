package store

import "github.com/winterochek/todo-server/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	FindById(int) (*model.User, error)
	FindAll() ([]*model.User, error)
}

type TaskRepository interface {
	Create(*model.Task) error
}
