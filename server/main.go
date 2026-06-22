package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/van5h1007/GoToDo/internal/handler"
	"github.com/van5h1007/GoToDo/internal/repository"
	"github.com/van5h1007/GoToDo/internal/service"
)

func main() {
	repo := repository.NewInMemoryTodoRepository()
	service := service.NewTodoService(repo)
	todoHandler := handler.NewTodoHandler(service)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	api := router.Group("/api/v1")
	todoHandler.RegisterRoutes(api)

	log.Println("Server running on port:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
