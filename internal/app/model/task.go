package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Task struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Content, validation.Required, validation.Length(3, 159)),
	)
}
