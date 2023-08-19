package api

import (
	"encoding/json"
	"net/http"

	"github.com/winterochek/todo-server/internal/app/model"
	h "github.com/winterochek/todo-server/internal/helpers"
)

// users/create POST controller
func (s *server) HandleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	type response struct {
		User  *model.User `json:"user"`
		Token string      `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := model.User{
			Email:    req.Email,
			Password: req.Password,
			Username: req.Username,
		}

		if err := u.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := u.BeforeCreate(); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.store.User().Create(&u)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u.Sanitaze()
		token, err := s.jwtClient.GenerateToken(u.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}

		res := &response{
			User:  &u,
			Token: token,
		}
		s.respond(w, r, http.StatusCreated, res)
	}
}

// users POST controller
func (s *server) HandleUsersLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User  *model.User `json:"user"`
		Token string      `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePasswords(req.Password) {
			s.error(w, r, http.StatusUnauthorized, ErrIncorrectEmailOrPassword)
			return
		}

		token, err := s.jwtClient.GenerateToken(u.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}
		res := &response{
			User:  u,
			Token: token,
		}
		s.respond(w, r, http.StatusOK, res)
	}
}

// users GET controller
func (s *server) HandleGetUsers() http.HandlerFunc {
	type response struct {
		Users  []*model.User `json:"users"`
		UserId int           `json:"userId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := h.GetUsedIDFromContext(r.Context(), context_userId_key)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
			return
		}

		users, err := s.store.User().FindAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
		}

		res := &response{
			Users:  users,
			UserId: userId,
		}

		s.respond(w, r, http.StatusOK, res)
	}
}