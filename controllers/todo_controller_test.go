package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"tuhuynh.com/go-ioc-gin-example/entities"
	"tuhuynh.com/go-ioc-gin-example/security"
)

// MockTodoService is a mock implementation of TodoService
type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) List(ctx context.Context) ([]entities.Todo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Todo), args.Error(1)
}

func (m *MockTodoService) Create(ctx context.Context, todo entities.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoService) Get(ctx context.Context, id int) (entities.Todo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entities.Todo), args.Error(1)
}

func (m *MockTodoService) Update(ctx context.Context, todo entities.Todo) error {
	args := m.Called(ctx, todo)
	return args.Error(0)
}

func (m *MockTodoService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTest() (*gin.Engine, *MockTodoService) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockService := new(MockTodoService)

	rateLimiter := &security.RateLimiter{}
	rateLimiter.PostConstruct()

	controller := &TodoController{
		Service:     mockService,
		RateLimiter: rateLimiter,
	}

	r.GET("/todos", controller.ListTodos)
	r.POST("/todos", controller.CreateTodo)
	r.GET("/todos/:id", controller.GetTodo)
	r.PUT("/todos/:id", controller.UpdateTodo)
	r.DELETE("/todos/:id", controller.DeleteTodo)

	return r, mockService
}

func TestListTodos(t *testing.T) {
	r, mockService := setupTest()

	t.Run("success", func(t *testing.T) {
		todos := []entities.Todo{{ID: 1, Title: "Test Todo", Completed: false}}
		mockService.On("List", mock.Anything).Return(todos, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/todos", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []entities.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, todos, response)
	})

	t.Run("error", func(t *testing.T) {
		mockService.On("List", mock.Anything).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/todos", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreateTodo(t *testing.T) {
	r, mockService := setupTest()

	t.Run("success", func(t *testing.T) {
		todo := entities.Todo{Title: "New Todo", Completed: false}
		mockService.On("Create", mock.Anything, todo).Return(nil).Once()

		todoJSON, _ := json.Marshal(todo)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(todoJSON))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetTodo(t *testing.T) {
	r, mockService := setupTest()

	t.Run("success", func(t *testing.T) {
		todo := entities.Todo{ID: 1, Title: "Test Todo", Completed: false}
		mockService.On("Get", mock.Anything, 1).Return(todo, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/todos/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response entities.Todo
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, todo, response)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("Get", mock.Anything, 999).Return(entities.Todo{}, errors.New("not found")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/todos/999", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdateTodo(t *testing.T) {
	r, mockService := setupTest()

	t.Run("success", func(t *testing.T) {
		todo := entities.Todo{ID: 1, Title: "Updated Todo", Completed: true}
		mockService.On("Update", mock.Anything, todo).Return(nil).Once()

		todoJSON, _ := json.Marshal(todo)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/todos/1", bytes.NewBuffer(todoJSON))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/todos/1", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteTodo(t *testing.T) {
	r, mockService := setupTest()

	t.Run("success", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, 1).Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/todos/1", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/todos/invalid", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
