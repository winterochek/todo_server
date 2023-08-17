package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	jwtclient "github.com/winterochek/todo-server/internal/app/jwt-client"
	"github.com/winterochek/todo-server/internal/app/model"
	"github.com/winterochek/todo-server/internal/app/store"
)

var (
	ErrIncorrectEmailOrPassword = errors.New("incorrent email or password")
	ErrInternal                 = errors.New("internal server error")
)

type server struct {
	router    *mux.Router
	store     store.Store
	jwtClient *jwtclient.JWTClient
}

func NewServer(st store.Store, jwtClient *jwtclient.JWTClient) *server {
	s := &server{
		router:    mux.NewRouter(),
		store:     st,
		jwtClient: jwtClient,
	}
	s.ConfigureRouter()
	return s
}

func (s *server) ConfigureRouter() {
	s.router.HandleFunc("/users/create", s.HandleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/users", s.HandleUsersLogin()).Methods("POST")
	s.router.HandleFunc("/users", s.AuthMiddleware(s.HandleGetUsers())).Methods("GET")
	s.router.HandleFunc("/tasks", s.AuthMiddleware(s.HandleTaskCreate())).Methods("POST")
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

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

func (s *server) HandleGetUsers() http.HandlerFunc {
	type response struct {
		Users  []*model.User `json:"users"`
		UserId int           `json:"userId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value(context_userId_key).(int)
		if !ok {
			s.error(w, r, http.StatusInternalServerError, ErrInternal)
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
