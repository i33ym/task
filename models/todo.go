package models

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrTaskNotFound = errors.New("no task has been found")
)

type ToDo struct {
	ID          int
	CreatedAt   time.Time
	Title       string
	Description string
	Deadline    time.Time
	Priority    string
	Done        bool
	UpdatedAt   time.Time
	Version     int
}

func (models *Models) Create(todo *ToDo) error {
	stmt := `insert into todo (title, description, priority, deadline, done, updated_at)
	values (?, ?, ?, ?, ?, ?);`

	if _, err := models.db.Exec(stmt, todo.Title, todo.Description, todo.Priority, todo.Deadline, todo.Done, todo.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (models *Models) Update(todo *ToDo) error {
	return nil
}

func (models *Models) Fetch(id int) (*ToDo, error) {
	if id < 1 {
		return nil, ErrTaskNotFound
	}

	stmt := `select * from todo where id = ?;`
	row := models.db.QueryRow(stmt, id)
	todo := &ToDo{}
	if err := row.Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.Title,
		&todo.Description,
		&todo.Priority,
		&todo.Deadline,
		&todo.Done,
		&todo.UpdatedAt,
		&todo.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrTaskNotFound
		default:
			return nil, err
		}
	}

	return todo, nil
}

func (models *Models) Delete(id int) error {
	return nil
}
