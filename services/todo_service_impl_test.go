package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"tuhuynh.com/go-ioc-gin-example/cache"
	"tuhuynh.com/go-ioc-gin-example/entities"
	"tuhuynh.com/go-ioc-gin-example/repositories"
)

func setupTestService() (*TodoServiceImpl, *repositories.TodoCrudRepositoryMock, *cache.RedisMock) {
	mockRepo := &repositories.TodoCrudRepositoryMock{}
	mockCache := &cache.RedisMock{}
	service := &TodoServiceImpl{
		Repository: mockRepo,
		Cache:      mockCache,
	}
	return service, mockRepo, mockCache
}

func TestTodoServiceImpl_Create(t *testing.T) {
	service, repo, cache := setupTestService()
	ctx := context.Background()

	todo := entities.Todo{
		Title:     "Test Todo",
		Completed: false,
	}

	err := service.Create(ctx, todo)
	assert.NoError(t, err)

	// Verify todo was created in repository
	todos, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, todo.Title, todos[0].Title)

	// Verify cache was invalidated
	cachedList, err := cache.Get(ctx, "todos:list")
	assert.NoError(t, err)
	assert.Nil(t, cachedList)
}

func TestTodoServiceImpl_Get(t *testing.T) {
	service, repo, _ := setupTestService()
	ctx := context.Background()

	// Create a todo first
	todo := entities.Todo{
		Title:     "Test Todo",
		Completed: false,
	}
	err := repo.Create(ctx, todo)
	assert.NoError(t, err)

	// Get the created todo
	todos, err := repo.List(ctx)
	assert.NoError(t, err)
	createdTodo := todos[0]

	// Test getting the todo
	fetchedTodo, err := service.Get(ctx, createdTodo.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdTodo.ID, fetchedTodo.ID)
	assert.Equal(t, createdTodo.Title, fetchedTodo.Title)
}

func TestTodoServiceImpl_List(t *testing.T) {
	service, repo, _ := setupTestService()
	ctx := context.Background()

	// Create some test todos
	todos := []entities.Todo{
		{Title: "Todo 1"},
		{Title: "Todo 2"},
	}

	for _, todo := range todos {
		err := repo.Create(ctx, todo)
		assert.NoError(t, err)
	}

	// Test listing todos
	fetchedTodos, err := service.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, fetchedTodos, 2)
}

func TestTodoServiceImpl_Update(t *testing.T) {
	service, repo, cache := setupTestService()
	ctx := context.Background()

	// Create a todo first
	todo := entities.Todo{
		Title:     "Test Todo",
		Completed: false,
	}
	err := repo.Create(ctx, todo)
	assert.NoError(t, err)

	todos, err := repo.List(ctx)
	assert.NoError(t, err)
	createdTodo := todos[0]

	// Update the todo
	createdTodo.Title = "Updated Title"
	createdTodo.Completed = true

	err = service.Update(ctx, createdTodo)
	assert.NoError(t, err)

	// Verify update
	updatedTodo, err := service.Get(ctx, createdTodo.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedTodo.Title)
	assert.True(t, updatedTodo.Completed)

	// Verify caches were invalidated
	cachedList, err := cache.Get(ctx, "todos:list")
	assert.NoError(t, err)
	assert.Nil(t, cachedList)
}

func TestTodoServiceImpl_Delete(t *testing.T) {
	service, repo, cache := setupTestService()
	ctx := context.Background()

	// Create a todo first
	todo := entities.Todo{
		Title: "Test Todo",
	}
	err := repo.Create(ctx, todo)
	assert.NoError(t, err)

	todos, err := repo.List(ctx)
	assert.NoError(t, err)
	createdTodo := todos[0]

	// Delete the todo
	err = service.Delete(ctx, createdTodo.ID)
	assert.NoError(t, err)

	// Verify deletion
	_, err = service.Get(ctx, createdTodo.ID)
	assert.Error(t, err)

	// Verify caches were invalidated
	cachedList, err := cache.Get(ctx, "todos:list")
	assert.NoError(t, err)
	assert.Nil(t, cachedList)

	cachedItem, err := cache.Get(ctx, fmt.Sprintf("todos:%d", createdTodo.ID))
	assert.NoError(t, err)
	assert.Nil(t, cachedItem)
}
