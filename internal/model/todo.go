package model

import "time"

type TodoStatus string

const (
	StatusPending    TodoStatus = "Pending"
	StatusInProgress TodoStatus = "In_Progress"
	StatusDone       TodoStatus = "Done"
)

type Todo struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TodoStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title string `json:"title" binding:"required, min=1, max=200"`
	Description string `json:"description" binding:"max=1000"`
}

type UpdateTodoRequest struct {
	Title *string `json:"title" binding:"omitempty, min=1, max=200"`
	Description *string `json:"description" binding:"omitempty, max=1000"`
	Status *TodoStatus `json:"status" binding:"omitempty"`
}
