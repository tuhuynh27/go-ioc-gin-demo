package migrations

import (
	"gorm.io/gorm"
	"tuhuynh.com/go-ioc-gin-example/config"
	"tuhuynh.com/go-ioc-gin-example/logger"
)

// Runner handles all database migrations
type Runner struct {
	Component struct{}
	Log       logger.Logger  `autowired:"true"`
	Config    *config.Config `autowired:"true"`
}

// Run executes all migrations
func (r *Runner) Run() error {
	r.Log.Info("Starting database migrations...")

	db := r.Config.DB

	// Add all migrations here
	migrations := []func(*gorm.DB) error{
		TodoMigration,
	}

	// Execute each migration
	for _, migration := range migrations {
		if err := migration(db); err != nil {
			return err
		}
	}

	r.Log.Info("All migrations completed successfully")
	return nil
}
