package teststore

import (
	"github.com/winterochek/todo-server/internal/app/model"
	"github.com/winterochek/todo-server/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (ur *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(ur.users) + 1
	ur.users[u.ID] = u

	return nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {

	for _, u := range ur.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.ErrRecordNotFound

}

func (ur *UserRepository) FindById(id int) (*model.User, error) {
	u, ok := ur.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (ur *UserRepository) FindAll() ([]*model.User, error) {
	users := []*model.User{}
	for _, v := range ur.users {
		users = append(users, v)
	}
	return users, nil
}
