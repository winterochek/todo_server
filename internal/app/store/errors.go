package store

import "errors"

var (
	ErrRecordNotFound         = errors.New("record not found")
	ErrEmailOrUsernameIsTaken = errors.New("email or username is taken")
)
