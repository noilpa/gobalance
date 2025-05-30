package server

import (
	"sync/atomic"
	"testing"
)

func TestSelectBackendRoundRobin(t *testing.T) {
	backends := []*backend{
		{address: "localhost:9001"},
		{address: "localhost:9002"},
		{address: "localhost:9003"},
	}

	s := &Server{
		strategy: "roundrobin",
		backends: backends,
	}

	// Прогоним 6 выборов — должен вернуться 0,1,2,0,1,2
	expectedOrder := []string{
		"localhost:9001",
		"localhost:9002",
		"localhost:9003",
		"localhost:9001",
		"localhost:9002",
		"localhost:9003",
	}

	for i, expected := range expectedOrder {
		b := s.selectBackend()
		if b.address != expected {
			t.Errorf("at step %d: expected %s, got %s", i, expected, b.address)
		}
	}

	// Проверка переполнения (симулируем counter близкий к max)
	atomic.StoreUint64(&s.counter, ^uint64(0)-1) // max-1
	b1 := s.selectBackend()
	b2 := s.selectBackend()
	if b1 == nil || b2 == nil {
		t.Fatal("backend is nil near counter overflow")
	}
	t.Logf("Successfully handled overflow: got %s and %s", b1.address, b2.address)
}

func TestSelectBackendLeastConn(t *testing.T) {
	// Создаём три фейковых бекенда
	b1 := &backend{address: "127.0.0.1:9001"}
	b2 := &backend{address: "127.0.0.1:9002"}
	b3 := &backend{address: "127.0.0.1:9003"}

	// Назначаем нагрузку
	atomic.StoreInt64(&b1.active, 5)
	atomic.StoreInt64(&b2.active, 2)
	atomic.StoreInt64(&b3.active, 10)

	s := &Server{
		strategy: "leastconn",
		backends: []*backend{b1, b2, b3},
	}

	selected := s.selectBackend()
	if selected != b2 {
		t.Errorf("Expected backend b2 (least active), got %v", selected.address)
	}
}
