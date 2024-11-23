package repositories

import (
	"context"

	"gorm.io/gorm"
	"tuhuynh.com/go-ioc-gin-example/config"
	"tuhuynh.com/go-ioc-gin-example/entities"
)

type TodoCrudRepositorySql struct {
	Component struct{}       `implements:"TodoCrudRepository"`
	Qualifier struct{}       `value:"sql"`
	Config    *config.Config `autowired:"true"`
}

func (r *TodoCrudRepositorySql) List(ctx context.Context) ([]entities.Todo, error) {
	var todos []entities.Todo
	result := r.Config.DB.WithContext(ctx).Find(&todos)
	return todos, result.Error
}

func (r *TodoCrudRepositorySql) Create(ctx context.Context, todo entities.Todo) error {
	return r.Config.DB.WithContext(ctx).Create(&todo).Error
}

func (r *TodoCrudRepositorySql) Get(ctx context.Context, id int) (entities.Todo, error) {
	var todo entities.Todo
	result := r.Config.DB.WithContext(ctx).First(&todo, id)
	return todo, result.Error
}

func (r *TodoCrudRepositorySql) Update(ctx context.Context, todo entities.Todo) error {
	result := r.Config.DB.WithContext(ctx).Save(&todo)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *TodoCrudRepositorySql) Delete(ctx context.Context, id int) error {
	result := r.Config.DB.WithContext(ctx).Delete(&entities.Todo{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
