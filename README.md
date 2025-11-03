# di - Simple Dependency Injection for Go

A lightweight, thread-safe, and lazy-initialized **Dependency Injection (DI)** library for Go.
Designed to make managing dependencies easy in any Go project, from small utilities to large applications.

## Features

- [x] **Lazy initialization** - instance is only created when first requested
- [x] **Singleton by default** - one instance per type
- [x] **Thread-safe** - safe for concurrent use
- [x] **Supports interfaces and structs** - inject dependencies of any type
- [x] **Remove/unregister instances** - free memory when needed

## Installation

```shell
go get -u github.com/ryanbekhen/di
```

## Usage

### Basi example

```go
package main

import (
	"github.com/ryanbekhen/di"
	"fmt"
)

// Service interface
type Service interface {
	Do()
}

type ServiceImpl struct{}

func (s *ServiceImpl) Do() {
	fmt.Println("Service running")
}

func main() {
	// Lazy initialization with factory
	di.RegisterFactory[Service](func() Service {
		fmt.Println("Creating Service now...")
		return &ServiceImpl{}
	})

	// Resolve instance
	s := di.MustResolve[Service]()
	s.Do()
}
```

### Example with dependencies

```go
package main

import (
	"github.com/ryanbekhen/di"
	"fmt"
)

// Service interface
type Service interface {
	Do()
}

type ServiceImpl struct{}

func (s *ServiceImpl) Do() {
	fmt.Println("Service running")
}

func main() {
	// Lazy initialization with factory
	di.RegisterFactory[Service](func() Service {
		fmt.Println("Creating Service now...")
		return &ServiceImpl{}
	})

	// Resolve instance
	s := di.MustResolve[Service]()
	s.Do()
}
```

### Example with dependencies

```go
package main

import (
	"github.com/ryanbekhen/di"
	"fmt"
)

// ------------------ Database Layer ------------------

// DBClient represents a simple database client
type DBClient struct{}

func (db *DBClient) Query() string {
	return "Query result"
}

// ------------------ Repository Layer ------------------

// UserRepository depends on DBClient
type UserRepository struct {
	DB *DBClient
}

// Factory function for UserRepository
func NewUserRepository() *UserRepository {
	db := di.MustResolve[*DBClient]()           // Resolve DBClient from DI
	return &UserRepository{DB: db}
}

// ------------------ Service Layer ------------------

// UserService depends on UserRepository
type UserService struct {
	Repo *UserRepository
}

// Factory function for UserService
func NewUserService() *UserService {
	repo := di.MustResolve[*UserRepository]()  // Resolve UserRepository from DI
	return &UserService{Repo: repo}
}

// ------------------ Main ------------------

func main() {
	// Register DB client directly
	di.Register[*DBClient](&DBClient{})

	// Register factories for dependent objects
	di.RegisterFactory[*UserRepository](NewUserRepository)
	di.RegisterFactory[*UserService](NewUserService)

	// Resolve service
	service := di.MustResolve[*UserService]()
	fmt.Println(service.Repo.DB.Query()) // Output: Query result
}
```

## API

See [API documentation](https://pkg.go.dev/github.com/ryanbekhen/di)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

MIT License - free to use, modify, and distribute