package main

import (
    "io"
    "log"
    "net"
    "sync/atomic"
)

var backends = []string{"localhost:9001", "localhost:9002"}
var counter uint64 = 0

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    log.Println("Listening on :8080")

    for {
        clientConn, err := listener.Accept()
        if err != nil {
            log.Println("Failed to accept:", err)
            continue
        }

        backend := selectBackend()
        go handleConnection(clientConn, backend)
    }
}

func selectBackend() string {
    idx := atomic.AddUint64(&counter, 1)
    return backends[int(idx-1)%len(backends)]
}

func handleConnection(clientConn net.Conn, backend string) {
    defer clientConn.Close()

    backendConn, err := net.Dial("tcp", backend)
    if err != nil {
        log.Println("Failed to connect to backend:", err)
        return
    }
    defer backendConn.Close()

    // Проксируем данные в обе стороны
    go io.Copy(backendConn, clientConn)
    io.Copy(clientConn, backendConn)
}
