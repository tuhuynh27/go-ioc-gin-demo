package services

import (
	"context"

	"tuhuynh.com/go-ioc-gin-example/entities"
)

type TodoService interface {
	List(ctx context.Context) ([]entities.Todo, error)
	Create(ctx context.Context, todo entities.Todo) error
	Get(ctx context.Context, id int) (entities.Todo, error)
	Update(ctx context.Context, todo entities.Todo) error
	Delete(ctx context.Context, id int) error
}
