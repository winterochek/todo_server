package api

import (
	"database/sql"
	"net/http"

	jwtclient "github.com/winterochek/todo-server/internal/app/jwt-client"
	sqlstore "github.com/winterochek/todo-server/internal/app/store/sql-store"
	"github.com/winterochek/todo-server/migrations"
)

// TOD0: move sensitive information to environment variables

func Start() error {
	db, err := NewDB("host=localhost port=5432 user=admin dbname=postgres password=admin sslmode=disable")
	if err != nil {
		return err
	}

	migrations.Up(db)
	defer db.Close()

	st := sqlstore.New(db)
	jwtClient := jwtclient.New("secret", 24)

	srv := NewServer(st, jwtClient)
	return http.ListenAndServe(":8000", srv)
}

func NewDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
