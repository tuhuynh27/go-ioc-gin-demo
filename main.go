package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"tuhuynh.com/go-ioc-gin-example/wire"
)

func main() {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found: %v", err)
		}
	}

	container, cleanup := wire.Initialize()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run application in goroutine
	go func() {
		container.Application.Run()
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down gracefully...")

	cleanup()
}
