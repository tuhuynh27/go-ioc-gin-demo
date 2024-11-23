package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Component struct{}
	DB        *gorm.DB
	Redis     *redis.Client
	Port      string
	AppMode   string
}

func NewConfig() *Config {
	appPort := getEnvOrDefault("APP_PORT", "8080")
	appMode := getEnvOrDefault("APP_MODE", "local")

	return &Config{
		DB:      initDB(),
		Redis:   initRedis(),
		Port:    fmt.Sprintf(":%s", appPort),
		AppMode: appMode,
	}
}

func (c *Config) PreDestroy() {
	// Close Redis connection
	if c.Redis != nil {
		if err := c.Redis.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		} else {
			log.Println("Successfully closed Redis connection")
		}
	}

	// Close DB connection
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			log.Printf("Error getting underlying *sql.DB: %v", err)
		} else {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error closing DB connection: %v", err)
			} else {
				log.Println("Successfully closed DB connection")
			}
		}
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initDB() *gorm.DB {
	dbUser := getEnvOrDefault("DB_USER", "myuser")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "mypassword")
	dbName := getEnvOrDefault("DB_NAME", "mydatabase")
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "3306")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying *sql.DB: %v", err)
	}

	// Test the connection
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db
}

func initRedis() *redis.Client {
	log.Println("Initializing redis client")
	redisHost := getEnvOrDefault("REDIS_HOST", "localhost")
	redisPort := getEnvOrDefault("REDIS_PORT", "6379")
	redisPassword := getEnvOrDefault("REDIS_PASSWORD", "")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0, // use default DB
	})

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("failed to ping redis: %v", err)
		return nil
	}

	log.Println("Successfully connected to redis")
	return client
}
