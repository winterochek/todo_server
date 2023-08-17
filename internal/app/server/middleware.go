package api

import (
	"context"
	"errors"
	"net/http"

	h "github.com/winterochek/todo-server/internal/helpers"
)

const (
	header_authorization_key = "Authorization"
	context_userId_key       = "userId"
)

var (
	ErrNoAuthorizationHeader          = errors.New("Authorization header is missing")
	ErrWrongTypeOfAuthorizationHeader = errors.New("Authorization header type is wrong")
	ErrNotAuthorized                  = errors.New("Not authorized")
	ErrNotFound                       = errors.New("Not found")
)

func (s *server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(header_authorization_key)
		tokenString, err := h.SpliceToken(authHeader)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, ErrNotAuthorized)
			return
		}
		id, err := s.jwtClient.ParseToken(tokenString)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, ErrNotAuthorized)
			return
		}

		ctx := context.WithValue(r.Context(), context_userId_key, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
