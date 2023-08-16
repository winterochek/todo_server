package store

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEmailIsTaken   = errors.New("email is taken")
)