package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	jwtclient "github.com/winterochek/todo-server/internal/app/jwt-client"
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
	s.router.HandleFunc("/tasks", s.AuthMiddleware(s.HandleMultipleTasksRead())).Methods("GET")
	s.router.HandleFunc("/tasks/{id:[0-9]+}", s.AuthMiddleware(s.HandleSingleTaskRead())).Methods("GET")
	s.router.HandleFunc("/tasks/{id:[0-9]+}", s.AuthMiddleware(s.HandleTasksUpdate())).Methods("PUT")
	s.router.HandleFunc("/tasks/{id:[0-9]+}", s.AuthMiddleware(s.HandleTasksDelete())).Methods("DELETE")
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
