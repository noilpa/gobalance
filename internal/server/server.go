package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type backend struct {
	address string
	active  int64 // атомарный счетчик активных соединений
}

type Server struct {
	port     int
	strategy string
	backends []*backend
	counter  uint64
}

func New(port int, strategy string, backendAddrs []string) (*Server, error) {
	backends := make([]*backend, len(backendAddrs))
	for i, addr := range backendAddrs {
		backends[i] = &backend{address: addr}
	}
	return &Server{
		port:     port,
		strategy: strategy,
		backends: backends,
	}, nil
}

func (s *Server) Start() error {
	address := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Listening on " + address)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept:", err)
			continue
		}

		backend := s.selectBackend()
		go handleConnection(clientConn, backend)
	}
}

func (s *Server) selectBackend() *backend {
	switch s.strategy {
	case "roundrobin":
		idx := atomic.AddUint64(&s.counter, 1)
		id := int((idx - 1) % uint64(len(s.backends)))
		return s.backends[id]
	case "leastconn":
		var min *backend
		for _, b := range s.backends {
			if min == nil || atomic.LoadInt64(&b.active) < atomic.LoadInt64(&min.active) {
				min = b
			}
		}
		return min
	default:
		log.Printf("Unknown strategy: %s, fallback to roundrobin", s.strategy)
		idx := atomic.AddUint64(&s.counter, 1)
		return s.backends[int(idx-1)%len(s.backends)]
	}
}

func handleConnection(clientConn net.Conn, b *backend) {
	defer clientConn.Close()

	atomic.AddInt64(&b.active, 1)
	defer atomic.AddInt64(&b.active, -1)

	backendConn, err := net.Dial("tcp", b.address)
	if err != nil {
		log.Println("Failed to connect to backend:", err)
		return
	}
	defer backendConn.Close()

	go io.Copy(backendConn, clientConn)
	io.Copy(clientConn, backendConn)
}
