package di

import (
	"fmt"
	"reflect"
	"sync"
)

// instances stores singleton instances
var instances sync.Map

// factories stores factory functions for lazy initialization
var factories sync.Map

// typeKey returns a unique string key for any generic type (including interfaces)
func typeKey[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}

// Register registers a singleton instance directly
func Register[T any](instance T) {
	key := typeKey[T]()
	instances.Store(key, instance)
}

// RegisterFactory registers a factory function for lazy initialization
func RegisterFactory[T any](f func() T) {
	key := typeKey[T]()
	factories.Store(key, f)
}

// Resolve retrieves an instance from the container
func Resolve[T any]() (T, error) {
	key := typeKey[T]()

	if v, ok := instances.Load(key); ok {
		return v.(T), nil
	}

	if f, ok := factories.Load(key); ok {
		instance := f.(func() T)()
		instances.Store(key, instance)
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
	key := typeKey[T]()
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
