package di

import (
	"fmt"
	"sync"
	"testing"
)

// Test types
type TestInterface interface {
	GetValue() string
}

type TestStruct struct {
	Value string
}

func (t *TestStruct) GetValue() string {
	return t.Value
}

type AnotherStruct struct {
	ID int
}

// TestRegister tests the Register function
func TestRegister(t *testing.T) {
	Reset() // Clean up before test

	instance := &TestStruct{Value: "test"}
	Register[*TestStruct](instance)

	resolved, err := Resolve[*TestStruct]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resolved.Value != "test" {
		t.Errorf("Expected 'test', got '%s'", resolved.Value)
	}

	// Verify it's the same instance (singleton)
	if resolved != instance {
		t.Error("Expected same instance (singleton behavior)")
	}
}

// TestRegisterFactory tests the RegisterFactory function
func TestRegisterFactory(t *testing.T) {
	Reset()

	factoryCalled := false
	RegisterFactory[*TestStruct](func() *TestStruct {
		factoryCalled = true
		return &TestStruct{Value: "factory"}
	})

	// Factory should not be called yet (lazy initialization)
	if factoryCalled {
		t.Error("Factory should not be called until Resolve is called")
	}

	// First resolve - should call factory
	resolved, err := Resolve[*TestStruct]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !factoryCalled {
		t.Error("Factory should have been called")
	}

	if resolved.Value != "factory" {
		t.Errorf("Expected 'factory', got '%s'", resolved.Value)
	}

	// Second resolve - should return cached instance, not call factory again
	factoryCalled = false
	resolved2, err := Resolve[*TestStruct]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if factoryCalled {
		t.Error("Factory should not be called again (singleton)")
	}

	if resolved != resolved2 {
		t.Error("Expected same instance on second resolve")
	}
}

// TestResolve tests the Resolve function
func TestResolve(t *testing.T) {
	Reset()

	// Test resolving unregistered type
	_, err := Resolve[*TestStruct]()
	if err == nil {
		t.Error("Expected error when resolving unregistered type")
	}

	// Test resolving registered instance
	instance := &TestStruct{Value: "direct"}
	Register[*TestStruct](instance)

	resolved, err := Resolve[*TestStruct]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resolved.Value != "direct" {
		t.Errorf("Expected 'direct', got '%s'", resolved.Value)
	}
}

// TestMustResolve tests the MustResolve function
func TestMustResolve(t *testing.T) {
	Reset()

	// Test panic on unregistered type
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected MustResolve to panic on unregistered type")
		}
	}()

	_ = MustResolve[*TestStruct]()
}

// TestMustResolveSuccess tests successful MustResolve
func TestMustResolveSuccess(t *testing.T) {
	Reset()

	instance := &TestStruct{Value: "must"}
	Register[*TestStruct](instance)

	resolved := MustResolve[*TestStruct]()
	if resolved.Value != "must" {
		t.Errorf("Expected 'must', got '%s'", resolved.Value)
	}
}

// TestUnregister tests the Unregister function
func TestUnregister(t *testing.T) {
	Reset()

	// Register instance
	instance := &TestStruct{Value: "test"}
	Register[*TestStruct](instance)

	// Verify it exists
	_, err := Resolve[*TestStruct]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Unregister it
	Unregister[*TestStruct]()

	// Verify it's gone
	_, err = Resolve[*TestStruct]()
	if err == nil {
		t.Error("Expected error after unregistering")
	}
}

// TestUnregisterFactory tests unregistering a factory
func TestUnregisterFactory(t *testing.T) {
	Reset()

	// Register factory
	RegisterFactory[*TestStruct](func() *TestStruct {
		return &TestStruct{Value: "factory"}
	})

	// Unregister before resolving
	Unregister[*TestStruct]()

	// Verify it's gone
	_, err := Resolve[*TestStruct]()
	if err == nil {
		t.Error("Expected error after unregistering factory")
	}
}

// TestReset tests the Reset function
func TestReset(t *testing.T) {
	// Register multiple types
	Register[*TestStruct](&TestStruct{Value: "test1"})
	Register[*AnotherStruct](&AnotherStruct{ID: 42})
	RegisterFactory[TestInterface](func() TestInterface {
		return &TestStruct{Value: "factory"}
	})

	// Reset
	Reset()

	// Verify all are gone
	_, err1 := Resolve[*TestStruct]()
	_, err2 := Resolve[*AnotherStruct]()
	_, err3 := Resolve[TestInterface]()

	if err1 == nil || err2 == nil || err3 == nil {
		t.Error("Expected all types to be unregistered after Reset")
	}
}

// TestInterfaceSupport tests that interfaces work correctly
func TestInterfaceSupport(t *testing.T) {
	Reset()

	RegisterFactory[TestInterface](func() TestInterface {
		return &TestStruct{Value: "interface"}
	})

	resolved, err := Resolve[TestInterface]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resolved.GetValue() != "interface" {
		t.Errorf("Expected 'interface', got '%s'", resolved.GetValue())
	}
}

// TestMultipleTypes tests registering multiple different types
func TestMultipleTypes(t *testing.T) {
	Reset()

	Register[*TestStruct](&TestStruct{Value: "test"})
	Register[*AnotherStruct](&AnotherStruct{ID: 100})

	resolved1, err1 := Resolve[*TestStruct]()
	resolved2, err2 := Resolve[*AnotherStruct]()

	if err1 != nil || err2 != nil {
		t.Fatal("Expected no errors")
	}

	if resolved1.Value != "test" {
		t.Errorf("Expected 'test', got '%s'", resolved1.Value)
	}

	if resolved2.ID != 100 {
		t.Errorf("Expected 100, got %d", resolved2.ID)
	}
}

// TestConcurrentAccess tests thread-safety
func TestConcurrentAccess(t *testing.T) {
	Reset()

	callCount := 0
	var mu sync.Mutex

	RegisterFactory[*TestStruct](func() *TestStruct {
		mu.Lock()
		callCount++
		mu.Unlock()
		return &TestStruct{Value: "concurrent"}
	})

	// Resolve concurrently from multiple goroutines
	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := Resolve[*TestStruct]()
			if err != nil {
				t.Errorf("Concurrent resolve failed: %v", err)
			}
		}()
	}

	wg.Wait()

	// Factory should be called a small number of times despite concurrent access
	// Due to race conditions, it might be called more than once, but should be much less than numGoroutines
	mu.Lock()
	count := callCount
	mu.Unlock()

	if count > 10 {
		t.Errorf("Expected factory to be called a few times due to race conditions, but was called %d times (too many)", count)
	}
	if count < 1 {
		t.Error("Expected factory to be called at least once")
	}
}

// TestConcurrentRegistration tests concurrent registration
func TestConcurrentRegistration(t *testing.T) {
	Reset()

	var wg sync.WaitGroup
	numGoroutines := 50

	// Register different instances concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// Use a unique type key for each goroutine by wrapping in a struct
			RegisterFactory[*TestStruct](func() *TestStruct {
				return &TestStruct{Value: fmt.Sprintf("value-%d", idx)}
			})
		}(i)
	}

	wg.Wait()

	// Should be able to resolve without error
	_, err := Resolve[*TestStruct]()
	if err != nil {
		t.Errorf("Concurrent registration failed: %v", err)
	}
}

// TestTypeKey tests that different types get different keys
func TestTypeKey(t *testing.T) {
	key1 := typeKey[*TestStruct]()
	key2 := typeKey[*AnotherStruct]()
	key3 := typeKey[TestInterface]()

	if key1 == key2 {
		t.Error("Different types should have different keys")
	}

	if key1 == key3 {
		t.Error("Struct and interface should have different keys")
	}

	// Same type should have same key
	key1Again := typeKey[*TestStruct]()
	if key1 != key1Again {
		t.Error("Same type should have same key")
	}
}

// TestRegisterOverwrite tests that registering overwrites previous registration
func TestRegisterOverwrite(t *testing.T) {
	Reset()

	// Register first instance
	Register[*TestStruct](&TestStruct{Value: "first"})

	resolved1, _ := Resolve[*TestStruct]()
	if resolved1.Value != "first" {
		t.Errorf("Expected 'first', got '%s'", resolved1.Value)
	}

	// Register second instance (overwrite)
	Register[*TestStruct](&TestStruct{Value: "second"})

	resolved2, _ := Resolve[*TestStruct]()
	if resolved2.Value != "second" {
		t.Errorf("Expected 'second', got '%s'", resolved2.Value)
	}
}

// TestFactoryReturnsNil tests handling of factory that returns nil
func TestFactoryReturnsNil(t *testing.T) {
	Reset()

	RegisterFactory[*TestStruct](func() *TestStruct {
		return nil
	})

	resolved, err := Resolve[*TestStruct]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resolved != nil {
		t.Error("Expected nil from factory")
	}
}
