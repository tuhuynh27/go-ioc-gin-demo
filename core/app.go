package core

import (
	"tuhuynh.com/go-ioc-gin-example/logger"

	"github.com/gin-gonic/gin"
	"tuhuynh.com/go-ioc-gin-example/config"
	"tuhuynh.com/go-ioc-gin-example/controllers"
	"tuhuynh.com/go-ioc-gin-example/migrations"
)

type Application struct {
	Component       struct{}
	Config          *config.Config              `autowired:"true"`
	Log             logger.Logger               `autowired:"true"`
	HealthCheck     *HealthCheck                `autowired:"true"`
	TodoController  *controllers.TodoController `autowired:"true"`
	MigrationRunner *migrations.Runner          `autowired:"true"`
}

func (a *Application) Run() {
	// Run database migrations
	if err := a.MigrationRunner.Run(); err != nil {
		a.Log.Fatal("Failed to run migrations: %v", err)
	}

	router := gin.Default()

	// Health check endpoint
	router.GET("/health", a.HealthCheck.Check)

	router.GET("/todos", a.TodoController.ListTodos)
	router.POST("/todos", a.TodoController.CreateTodo)
	router.GET("/todos/:id", a.TodoController.GetTodo)
	router.PUT("/todos/:id", a.TodoController.UpdateTodo)
	router.DELETE("/todos/:id", a.TodoController.DeleteTodo)

	err := router.Run(a.Config.Port)
	if err != nil {
		a.Log.Fatal("Failed to start server: %v", err)
	}
}
