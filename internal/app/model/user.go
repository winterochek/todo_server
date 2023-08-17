package model

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	Username          string    `json:"username"`
	EncryptedPassword string    `json:"-"`
	Password          string    `json:"password,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
}

// Hash password before saving to DB
func (u *User) BeforeCreate() error {

	enc, err := encryptPassword(u.Password)
	if err != nil {
		return err
	}

	u.EncryptedPassword = enc
	return nil
}

// Validation of user model
func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(
			&u.Password,
			// plain password is required in case no encrypted pass provided
			validation.By(validation.RuleFunc(requiredIf(u.EncryptedPassword == ""))),
			validation.Length(8, 20)),
		validation.Field(
			&u.Username,
			validation.Required,
			validation.Length(3, 20),
		),
	)
}

// Remove unwanted fields
func (u *User) Sanitaze() {
	u.Password = ""
}

// Compare provided password with the enctypted one
func (u *User) ComparePasswords(p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(p)) == nil
}

// Enctypt password with bcrypt lib
func encryptPassword(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
