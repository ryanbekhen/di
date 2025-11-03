package di

import (
	"fmt"
	"sync"
)

// instances stores singleton instances
var instances sync.Map

// factories stores factory functions for lazy initialization
var factories sync.Map

// Register registers a singleton instance directly
func Register[T any](instance T) {
	key := fmt.Sprintf("%T", *new(T)) // type name as key
	instances.Store(key, instance)
}

// RegisterFactory registers a factory function for lazy initialization
// The instance is created only when Resolve is called for the first time
func RegisterFactory[T any](f func() T) {
	key := fmt.Sprintf("%T", *new(T))
	factories.Store(key, f)
}

// Resolve retrieves an instance from the container
// If the instance is registered via factory, it will be created lazily
func Resolve[T any]() (T, error) {
	key := fmt.Sprintf("%T", *new(T))

	// Check if instance already exists
	if v, ok := instances.Load(key); ok {
		return v.(T), nil
	}

	// Check if factory exists
	if f, ok := factories.Load(key); ok {
		instance := f.(func() T)()     // call factory function
		instances.Store(key, instance) // store for next calls
		return instance, nil
	}

	var zero T
	return zero, fmt.Errorf("no instance found for type %v", key)
}

// MustResolve retrieves an instance or panics if not found
func MustResolve[T any]() T {
	v, err := Resolve[T]()
	if err != nil {
		panic(err)
	}
	return v
}

// Unregister removes an instance or factory from the container
func Unregister[T any]() {
	key := fmt.Sprintf("%T", *new(T))
	instances.Delete(key)
	factories.Delete(key)
}

// Reset clears all instances and factories (useful for testing)
func Reset() {
	instances.Range(func(k, v any) bool {
		instances.Delete(k)
		return true
	})
	factories.Range(func(k, v any) bool {
		factories.Delete(k)
		return true
	})
}
