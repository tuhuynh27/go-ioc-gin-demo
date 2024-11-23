package repositories

import (
	"context"
	"database/sql"
	"sync"

	"tuhuynh.com/go-ioc-gin-example/entities"
)

type TodoCrudRepositoryMock struct {
	Component struct{} `implements:"TodoCrudRepository"`
	Qualifier struct{} `value:"mock"`
	todos     map[int]entities.Todo
	mutex     sync.RWMutex
	lastID    int
}

// Initialize the mock repository with an empty map
func (r *TodoCrudRepositoryMock) init() {
	if r.todos == nil {
		r.todos = make(map[int]entities.Todo)
	}
}

func (r *TodoCrudRepositoryMock) List(ctx context.Context) ([]entities.Todo, error) {
	r.init()
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	todos := make([]entities.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoCrudRepositoryMock) Create(ctx context.Context, todo entities.Todo) error {
	r.init()
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.lastID++
	todo.ID = r.lastID
	r.todos[todo.ID] = todo
	return nil
}

func (r *TodoCrudRepositoryMock) Get(ctx context.Context, id int) (entities.Todo, error) {
	r.init()
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if todo, exists := r.todos[id]; exists {
		return todo, nil
	}
	return entities.Todo{}, sql.ErrNoRows
}

func (r *TodoCrudRepositoryMock) Update(ctx context.Context, todo entities.Todo) error {
	r.init()
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return sql.ErrNoRows
	}

	r.todos[todo.ID] = todo
	return nil
}

func (r *TodoCrudRepositoryMock) Delete(ctx context.Context, id int) error {
	r.init()
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.todos[id]; !exists {
		return sql.ErrNoRows
	}

	delete(r.todos, id)
	return nil
}
