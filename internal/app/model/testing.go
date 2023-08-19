package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Email:    "example@mail.com",
		Username: "example",
		Password: "password",
	}
}
