package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/van5h1007/GoToDo/internal/model"
	"github.com/van5h1007/GoToDo/internal/repository"
	"github.com/van5h1007/GoToDo/internal/service"
)

// to create uniform API responses
type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type TodoHandler struct {
	service service.TodoService
}

func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{service: svc}
}

func (h *TodoHandler) RegisterRoutes(rg *gin.RouterGroup) {
	todos := rg.Group("/todos")
	todos.POST("", h.CreateTodo)
	todos.GET("", h.GetAllTodos)
	todos.GET("/:id", h.GetTodoByID)
	todos.PUT("/:id", h.UpdateTodo)
	todos.DELETE("/:id", h.DeleteTodo)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req model.CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	todo, err := h.service.CreateTodo(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "Failed to create todo",
		})
		return
	}
	c.JSON(http.StatusCreated, APIResponse{
		Success: true, 
		Data: todo,
	})
}

func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	id := c.Param("id")

	todo, err := h.service.GetTodo(id)
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Error: "Todo not found.",
			})
			return 
		}
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: todo,
	})
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	todos, err := h.service.GetAllTodos()
	if err!= nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: "Internal error",
		})
		return
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: todos,
	})
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id:= c.Param("id")

	var req model.UpdateTodoRequest
	if err:= c.ShouldBindJSON(&req); err!=nil{
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error: err.Error(),
		})
		return
	}

	todo, err := h.service.UpdateTodo(id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Error: "Todo not found.",
			})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Success: false,
			Error: err.Error(),
		})
		return 
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: todo,
	})
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	if err:= h.service.DeleteTodo(id); err!= nil {
		if errors.Is(err, repository.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Error: "Todo not found.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error: "Internal error.",
		})
		return 
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: "Todo deleted.",
	})
}


