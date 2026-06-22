package service

import (
	"errors"

	"github.com/van5h1007/GoToDo/internal/model"
	"github.com/van5h1007/GoToDo/internal/repository"
)

type TodoService interface {
	CreateTodo(req *model.CreateTodoRequest) (*model.Todo, error)
	GetTodo(id string) (*model.Todo, error)
	GetAllTodos() ([]*model.Todo, error)
	UpdateTodo(id string, req *model.UpdateTodoRequest) (*model.Todo, error)
	DeleteTodo(id string) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(req *model.CreateTodoRequest) (*model.Todo, error) {
	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
	}
	return s.repo.Create(todo)
}

func (s *todoService) GetTodo(id string) (*model.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *todoService) GetAllTodos() ([]*model.Todo, error) {
	return s.repo.GetAll()
}

func (s *todoService) UpdateTodo(id string, req *model.UpdateTodoRequest) (*model.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Status != nil {
		if err := validateStatusTransaction(todo.Status, *req.Status); err != nil {
			return nil, err
		}
		todo.Status = *req.Status
	}
	return s.repo.Update(todo)
}

func (s *todoService) DeleteTodo(id string) error {
	return s.repo.Delete(id)
}

func validateStatusTransaction(from, to model.TodoStatus) error {
	allowed := map[model.TodoStatus][]model.TodoStatus{
		model.StatusPending:    {model.StatusInProgress},
		model.StatusInProgress: {model.StatusDone, model.StatusPending},
		model.StatusDone:       {},
	}

	for _, validNext := range allowed[from] {
		if validNext == to {
			return nil
		}
	}

	return errors.New("invalid status transaction: " + string(from) + "->" + string(to))
}
