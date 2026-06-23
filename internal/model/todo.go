//this is like a schema defining file

package model

import "time"

// enum sort of
type TodoStatus string

const (
	StatusPending    TodoStatus = "Pending"
	StatusInProgress TodoStatus = "In_Progress"
	StatusDone       TodoStatus = "Done"
)

type Todo struct {
	ID          string     `json:"id"          db:"id"`
	Title       string     `json:"title"       db:"title"`
	Description string     `json:"description" db:"description"`
	Status      TodoStatus `json:"status"      db:"status"`
	CreatedAt   time.Time  `json:"created_at"  db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"  db:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"max=1000"`
}

// * are used so that if only title was updated other things
// remain as it is and do not take default nil value
type UpdateTodoRequest struct {
	Title       *string     `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string     `json:"description" binding:"omitempty,max=1000"`
	Status      *TodoStatus `json:"status" binding:"omitempty"`
}
