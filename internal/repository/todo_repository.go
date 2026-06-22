//this tells us how data is stored and retrieved  

package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/van5h1007/GoToDo/internal/model"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
	ErrTodoExists   = errors.New("todo already exists")
)

//interface is like abstraction in oops 
//it says does not matter how you save the todo
//if it has these 5 methods it is todo repository

//Interface: Defines the rules.
//Struct: Holds the physical data (the map) and does the real work.
//Constructor: Packages the struct up nicely so it can be Injected elsewhere.

type TodoRepository interface {
	Create(todo *model.Todo) (*model.Todo, error)
	GetByID(id string) (*model.Todo, error)
	GetAll() ([]*model.Todo, error)
	Update(todo *model.Todo) (*model.Todo, error)
	Delete(id string) error
}

type InMemoryTodoRepository struct {
	mu    sync.RWMutex
	todos map[string]*model.Todo
}

func NewInMemoryTodoRepository() TodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[string]*model.Todo),
	}
}

func (r *InMemoryTodoRepository) Create(todo *model.Todo) (*model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo.ID = uuid.New().String()
	todo.CreatedAt = time.Now().UTC()
	todo.UpdatedAt = time.Now().UTC()
	todo.Status = model.StatusPending

	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *InMemoryTodoRepository) GetByID(id string) (*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrTodoNotFound
	}
	return todo, nil
}

func (r *InMemoryTodoRepository) GetAll() ([]*model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todos := make([]*model.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *InMemoryTodoRepository) Update(todo *model.Todo) (*model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return nil, ErrTodoNotFound
	}

	todo.UpdatedAt = time.Now().UTC()
	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *InMemoryTodoRepository) Delete(id string) (error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return ErrTodoNotFound
	}

	delete(r.todos, id)
	return nil
}
