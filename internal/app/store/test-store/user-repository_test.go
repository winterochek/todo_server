package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winterochek/todo-server/internal/app/model"
	"github.com/winterochek/todo-server/internal/app/store"
	teststore "github.com/winterochek/todo-server/internal/app/store/test-store"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	email := "example@mail.com"

	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email

	s.User().Create(u)
	u, err = s.User().FindByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByID(t *testing.T) {
	s := teststore.New()
	id := -10

	_, err := s.User().FindById(id)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)

	s.User().Create(u)
	u, err = s.User().FindById(u.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u)

}

func TestUserRepository_FindAll(t *testing.T) {
	s := teststore.New()
	users, err := s.User().FindAll()
	assert.Equal(t, users, []*model.User{})
	assert.Nil(t, err)
	u := model.TestUser(t)
	s.User().Create(u)
	users, err = s.User().FindAll()
	assert.Equal(t, users, []*model.User{u})
	assert.Nil(t, err)
}
