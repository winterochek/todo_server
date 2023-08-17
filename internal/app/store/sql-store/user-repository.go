package sqlstore

import (
	"database/sql"

	"github.com/winterochek/todo-server/internal/app/model"
	"github.com/winterochek/todo-server/internal/app/store"
)

type UserRepository struct {
	store *Store
}

// create new User
func (ur *UserRepository) Create(u *model.User) error {
	if err := ur.store.db.QueryRow(
		"INSERT INTO users (email, username, encrypted_password) VALUES ($1, $2, $3) RETURNING id, created_at",
		u.Email, u.Username, u.EncryptedPassword,
	).Scan(&u.ID, *&u.CreatedAt); err != nil {
		return store.ErrEmailOrUsernameIsTaken
	}
	return nil
}

// check for providing true errors for most cases, not the only email one
// Find user by email
func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := ur.store.db.QueryRow(
		"SELECT id, email, username, encrypted_password, created_at FROM users WHERE email = $1",
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.EncryptedPassword, &u.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// find user by id
func (ur *UserRepository) FindById(id int) (*model.User, error) {
	u := &model.User{}
	if err := ur.store.db.QueryRow(
		"SELECT id, email, username, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// find all users
func (ur *UserRepository) FindAll() ([]*model.User, error) {
	var users []*model.User

	rows, err := ur.store.db.Query("SELECT id, email, username FROM users")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := &model.User{}
		err := rows.Scan(&u.ID, &u.Email, &u.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
