package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/van5h1007/GoToDo/internal/model"
)

type PostgresTodoRepository struct {
	db *sqlx.DB
	//it sort of a connection pool managing multiple connections
}

func NewPostgresTodoRepositpry(db *sqlx.DB) TodoRepository {
	return &PostgresTodoRepository{db: db}
}

func (r *PostgresTodoRepository) Create(todo *model.Todo) (*model.Todo, error) {
	query := `
	    INSERT INTO todos (id, title, description, status, created_at, updated_at)
		VALUES (:id, :title, :description, :status, :created_at, :updated_at)
		RETURNING id, title, description, status, created_at, updated_at`

	todo.ID = uuid.New().String()
	todo.CreatedAt = time.Now().UTC()
	todo.UpdatedAt = time.Now().UTC()
	todo.Status = model.StatusPending

	rows, err := r.db.NamedQuery(query, todo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(todo); err != nil {
			return nil, err
		}
	}
	return todo, nil
}

func (r *PostgresTodoRepository) GetByID(id string) (*model.Todo, error) {
	var todo model.Todo
	query := `SELECT id, title, description, status, created_at, updated_at FROM todos WHERE id= $1`

	err := r.db.Get(&todo, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (r *PostgresTodoRepository) GetAll() ([]*model.Todo, error) {
	var todos []*model.Todo
	query := `SELECT id, title, description, status, created_at, updated_at FROM todos ORDER BY created_at DESC`

	if err := r.db.Select(&todos, query); err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *PostgresTodoRepository) Update(todo *model.Todo) (*model.Todo, error) {
	query := `
	    UPDATE todos
		SET title= :title, description= :description, status= :status, updated_at= :updated_at
		WHERE id= :id
		RETURNING id, title, description, status, created_at, updated_at`

	todo.UpdatedAt = time.Now().UTC()

	rows, err := r.db.NamedQuery(query, todo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(todo); err != nil {
			return nil, err
		}
		return todo, nil
	}
	return nil, ErrTodoNotFound
}

func (r *PostgresTodoRepository) Delete(id string) error {
	query := `DELETE FROM todos WHERE id= $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTodoNotFound
	}
	return nil
}
