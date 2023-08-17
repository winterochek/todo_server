package sqlstore

import (
	"time"

	"github.com/winterochek/todo-server/internal/app/model"
)

type TaskRepository struct {
	store *Store
}

func (tr *TaskRepository) Create(t *model.Task) error {
	query := "INSERT INTO tasks (user_id, content, completed, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at"
	err := tr.store.db.QueryRow(query, t.UserID, t.Content, t.Completed, time.Now(), time.Now()).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TaskRepository) ReadOne(id int, userId int) (*model.Task, error) {
	var t model.Task
	query := "SELECT content, completed, created_at, updated_at FROM tasks WHERE id = $1 AND user_id = $2"
	err := tr.store.db.QueryRow(query, id, userId).Scan(&t.Content, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	t.ID = id
	t.UserID = userId
	return &t, nil
}

func (tr *TaskRepository) ReadAll(userId int) ([]*model.Task, error) {
	var tasks []*model.Task
	query := "SELECT id, content, completed, created_at, updated_at FROM tasks  WHERE user_id =  $1"
	rows, err := tr.store.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		task := &model.Task{}
		task.UserID = userId
		err := rows.Scan(&task.ID, &task.Content, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (tr *TaskRepository) Update(t *model.Task) error {
	query := "UPDATE tasks SET content = $1, completed = $2, updated_at = $3 WHERE id = $4 AND user_id = $5 RETURNING updated_at, created_at"

	err := tr.store.db.QueryRow(query, t.Content, t.Completed, time.Now(), t.ID, t.UserID).Scan(&t.UpdatedAt, &t.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TaskRepository) Delete(id int, userId int) error {
	query := "DELETE FROM tasks WHERE id = $1 AND user_id = $2"
	_, err := tr.store.db.Exec(query, id, userId)
	if err != nil {
		return err
	}
	return nil
}
