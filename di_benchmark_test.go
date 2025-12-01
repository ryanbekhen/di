package di_test

import (
	"testing"

	"github.com/ryanbekhen/di"
)

type Service struct {
	Value int
}

func BenchmarkRegister(b *testing.B) {
	di.Reset()
	for i := 0; i < b.N; i++ {
		di.Register(Service{Value: i})
	}
}

func BenchmarkResolve(b *testing.B) {
	di.Reset()
	di.Register(Service{Value: 42})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[Service]()
	}
}

func BenchmarkMustResolve(b *testing.B) {
	di.Reset()
	di.Register(Service{Value: 42})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = di.MustResolve[Service]()
	}
}

func BenchmarkRegisterFactoryResolve(b *testing.B) {
	di.Reset()
	di.RegisterFactory(func() Service {
		return Service{Value: 100}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = di.Resolve[Service]()
	}
}

func BenchmarkUnregister(b *testing.B) {
	di.Reset()
	di.Register(Service{Value: 55})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		di.Unregister[Service]()
		di.Register(Service{Value: 55})
	}
}

// ------------------- Concurrent Benchmarks -------------------

func BenchmarkResolveConcurrent(b *testing.B) {
	di.Reset()
	di.Register(Service{Value: 42})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = di.Resolve[Service]()
		}
	})
}

func BenchmarkMustResolveConcurrent(b *testing.B) {
	di.Reset()
	di.Register(Service{Value: 42})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = di.MustResolve[Service]()
		}
	})
}

func BenchmarkRegisterFactoryResolveConcurrent(b *testing.B) {
	di.Reset()
	di.RegisterFactory(func() Service {
		return Service{Value: 100}
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = di.Resolve[Service]()
		}
	})
}

func BenchmarkRegisterUnregisterConcurrent(b *testing.B) {
	di.Reset()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			di.Register(Service{Value: 1})
			di.Unregister[Service]()
		}
	})
}
