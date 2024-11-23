package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tuhuynh.com/go-ioc-gin-example/entities"
	"tuhuynh.com/go-ioc-gin-example/security"
	"tuhuynh.com/go-ioc-gin-example/services"
)

type TodoController struct {
	Component   struct{}
	Service     services.TodoService  `autowired:"true"`
	RateLimiter *security.RateLimiter `autowired:"true"`
}

func (c *TodoController) ListTodos(ctx *gin.Context) {
	todos, err := c.Service.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todos)
}

func (c *TodoController) CreateTodo(ctx *gin.Context) {
	ip := ctx.ClientIP()

	if !c.RateLimiter.AllowRequest(ip) {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded. Try again later."})
		return
	}

	var todo entities.Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.Service.Create(ctx.Request.Context(), todo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Todo created successfully"})
}

func (c *TodoController) GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}
	todo, err := c.Service.Get(ctx.Request.Context(), idInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (c *TodoController) UpdateTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var todo entities.Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.Service.Update(ctx.Request.Context(), todo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
}

func (c *TodoController) DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}
	err = c.Service.Delete(ctx.Request.Context(), idInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
