package migrations

import "database/sql"

const (
	up   = "CREATE TABLE IF NOT EXISTS users (id bigserial not null primary key,  email varchar not null unique, encrypted_password varchar not null);"
	down = "DROP TABLE IF EXISTS users;"
)

func Up(db *sql.DB) error {
	_, err := db.Exec(up)
	if err != nil {
		return err
	}
	return nil
}
