package helpers

import (
	"errors"
	"strings"
)

const (
	NoAuthHeader        = "No Authorization header provided"
	WrongAuthHeaderType = "Authorization header type is wrong"
)

var (
	ErrNoAuthHeader        = errors.New(NoAuthHeader)
	ErrWrongAuthHeaderType = errors.New(WrongAuthHeaderType)
)

func SpliceToken(header string) (string, error) {
	if header == "" {
		return "", ErrNoAuthHeader
	}
	tokenString := strings.TrimPrefix(header, "Bearer ")
	if tokenString == "" {
		return "", ErrWrongAuthHeaderType
	}

	return tokenString, nil
}
