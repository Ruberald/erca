package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

func newBackend(urlStr string) *Backend {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	return &Backend{
		URL:          u,
		Alive:        true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(u),
	}
}

func (b *Backend) IsAlive() bool {
	return b.Alive
}

type ServerPool struct {
	Backends []*Backend
	current  uint64
}

func (s *ServerPool) Next() *Backend {
	// Round-robin balancing
	current := atomic.AddUint64(&s.current, 1)

	backend := s.Backends[(current-1)%uint64(len(s.Backends))]

	return backend
}

func startBackend(port string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from backend %s\n", port)
	})

	log.Printf("Backend running on port %s", port)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Start our fake backend servers.
	go startBackend("8081")
	go startBackend("8082")
	go startBackend("8083")

	serverPool := &ServerPool{
		Backends: []*Backend{
			newBackend("http://localhost:8081"),
			newBackend("http://localhost:8082"),
			newBackend("http://localhost:8083"),
		},
	}

	// Probably should refactor this outside
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		backend := serverPool.Next()

		fmt.Println(
			"Proxying request to backend:",
			backend.URL.String(),
		)

		backend.ReverseProxy.ServeHTTP(w, r)
	})

	log.Println("Load balancer running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
