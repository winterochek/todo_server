package sqlstore

import (
	"database/sql"
	"fmt"
	// "strings"
	"testing"

	"github.com/winterochek/todo-server/migrations"
)

func TestDB(t *testing.T, dbURI string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	migrations.Up(db)

	return db, func(tables ...string) {
		if len(tables) > 0 {
			for _, table := range tables {
				query := fmt.Sprintf("DELETE FROM %s", table)
				db.Exec(query)
			}
		}
		db.Close()
	}
}
