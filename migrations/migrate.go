package migrations

import "database/sql"

const (
	UsersUp = `	CREATE TABLE IF NOT EXISTS users (
			   	id bigserial not null primary key,  
				email varchar not null unique, 
				username varchar not null unique,
				encrypted_password varchar not null,
				created_at TIMESTAMP DEFAULT current_timestamp);
			`
	UsersDown = "DROP TABLE IF EXISTS users;"
	TasksUp   = `CREATE TABLE IF NOT EXISTS tasks (
				id bigserial not null primary key,
				user_id bigint not null,
				content text not null,
				completed boolean default false,
				created_at TIMESTAMP DEFAULT current_timestamp,
				updated_at TIMESTAMP DEFAULT current_timestamp,
				FOREIGN KEY (user_id) REFERENCES users (id)  
				)
				`
	TasksDown = "DROP TABLE IF EXISTS tasks"
)

func Up(db *sql.DB) error {
	_, err := db.Exec(UsersUp)
	if err != nil {
		return err
	}

	_, err = db.Exec(TasksUp)
	if err != nil {
		return err
	}
	return nil
}
