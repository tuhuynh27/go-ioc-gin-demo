# Go IoC Container Example with Gin

This project demonstrates dependency injection and inversion of control (IoC) patterns in Go using the Gin web framework. It implements a simple Todo API with Redis caching and MySQL database.

## Features

- REST API using Gin framework
- MySQL database with GORM
- Redis caching layer
- Structured logging with Zap
- Database migrations
- Environment configuration

## API Endpoints

- `GET /todos` - List all todos
- `POST /todos` - Create a new todo 
- `GET /todos/:id` - Get a specific todo
- `PUT /todos/:id` - Update a todo
- `DELETE /todos/:id` - Delete a todo

## Getting Started

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run `make deps` to install dependencies
4. Run `make migrate` to set up the database
5. Run `make run` to start the server

The API will be available at `http://localhost:8080` by default.

## Available Make Commands

- `make run` - Run the application
- `make test` - Run all tests
- `make test-coverage` - Run tests with coverage report
- `make build` - Build the application
- `make clean` - Clean build artifacts
- `make deps` - Install dependencies
- `make lint` - Run linter
- `make migrate` - Run database migrations
