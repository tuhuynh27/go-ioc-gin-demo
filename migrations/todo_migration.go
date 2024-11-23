package migrations

import (
	"gorm.io/gorm"
	"tuhuynh.com/go-ioc-gin-example/entities"
)

// TodoMigration handles the database schema for todos
func TodoMigration(db *gorm.DB) error {
	err := db.AutoMigrate(&entities.Todo{})
	if err != nil {
		return err
	}

	return nil
}
