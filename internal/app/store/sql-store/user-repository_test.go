package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winterochek/todo-server/internal/app/model"
	"github.com/winterochek/todo-server/internal/app/store"
	sqlstore "github.com/winterochek/todo-server/internal/app/store/sql-store"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, dbURI)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	u.Email = "email@email.com"
	u.BeforeCreate()
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, dbURI)
	defer teardown("users")

	s := sqlstore.New(db)
	email := "example@mail.com"

	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	s.User().Create(u)
	u, err = s.User().FindByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindById(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, dbURI)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)
	u, err := s.User().FindById(u.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, dbURI)
	defer teardown("users")

	s := sqlstore.New(db)
	users, err := s.User().FindAll()
	assert.NoError(t, err)
	assert.Equal(t, users, []*model.User{})
	assert.NotNil(t, users)

	u := model.TestUser(t)
	s.User().Create(u)
	u.Sanitaze()
	users, err = s.User().FindAll()
	assert.NoError(t, err)
	assert.Equal(t, users, []*model.User{u})
	assert.NotNil(t, users)
}
