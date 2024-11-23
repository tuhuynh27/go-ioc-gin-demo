package services

import (
	"context"
	"encoding/json"
	"fmt"

	"tuhuynh.com/go-ioc-gin-example/cache"
	"tuhuynh.com/go-ioc-gin-example/entities"
	"tuhuynh.com/go-ioc-gin-example/repositories"
)

type TodoServiceImpl struct {
	Component  struct{}                        `implements:"TodoService"`
	Repository repositories.TodoCrudRepository `autowired:"true" qualifier:"sql"`
	Cache      cache.Cache                     `autowired:"true" qualifier:"redis"`
}

func (s *TodoServiceImpl) List(ctx context.Context) ([]entities.Todo, error) {
	// Try to get from cache first
	if cached, err := s.Cache.Get(ctx, "todos:list"); err == nil && cached != nil {
		var todos []entities.Todo
		if err := json.Unmarshal([]byte(cached.(string)), &todos); err == nil {
			return todos, nil
		}
	}

	// If not in cache, get from repository
	todos, err := s.Repository.List(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the results
	if todosJson, err := json.Marshal(todos); err == nil {
		s.Cache.Set(ctx, "todos:list", string(todosJson))
	}

	return todos, nil
}

func (s *TodoServiceImpl) Create(ctx context.Context, todo entities.Todo) error {
	err := s.Repository.Create(ctx, todo)
	if err != nil {
		return err
	}

	// Invalidate list cache
	s.Cache.Set(ctx, "todos:list", nil)
	return nil
}

func (s *TodoServiceImpl) Get(ctx context.Context, id int) (entities.Todo, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("todos:%d", id)
	if cached, err := s.Cache.Get(ctx, cacheKey); err == nil && cached != nil {
		var todo entities.Todo
		if err := json.Unmarshal([]byte(cached.(string)), &todo); err == nil {
			return todo, nil
		}
	}

	// If not in cache, get from repository
	todo, err := s.Repository.Get(ctx, id)
	if err != nil {
		return todo, err
	}

	// Cache the result
	if todoJson, err := json.Marshal(todo); err == nil {
		s.Cache.Set(ctx, cacheKey, string(todoJson))
	}

	return todo, nil
}

func (s *TodoServiceImpl) Update(ctx context.Context, todo entities.Todo) error {
	err := s.Repository.Update(ctx, todo)
	if err != nil {
		return err
	}

	// Invalidate caches
	s.Cache.Set(ctx, fmt.Sprintf("todos:%d", todo.ID), nil)
	s.Cache.Set(ctx, "todos:list", nil)
	return nil
}

func (s *TodoServiceImpl) Delete(ctx context.Context, id int) error {
	err := s.Repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate caches
	s.Cache.Set(ctx, fmt.Sprintf("todos:%d", id), nil)
	s.Cache.Set(ctx, "todos:list", nil)
	return nil
}
