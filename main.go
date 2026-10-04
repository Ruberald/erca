package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy

	mux sync.RWMutex
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
	b.mux.RLock()
	defer b.mux.RUnlock()

	return b.Alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.Alive = alive
}

type ServerPool struct {
	Backends []*Backend
	current  uint64
}

func (s *ServerPool) Next() *Backend {
	current := atomic.AddUint64(&s.current, 1)

	for i := 0; i < len(s.Backends); i++ {
		index := int(current+uint64(i)) % len(s.Backends)

		backend := s.Backends[index]

		if backend.IsAlive() {
			return backend
		}
	}

	return nil
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	serverPool := &ServerPool{}

	for _, backendURL := range config.Backends {
		// Create the LB's representation of this backend.
		serverPool.Backends = append(
			serverPool.Backends,
			newBackend(backendURL),
		)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		backend := serverPool.Next()

		if backend == nil {
			http.Error(
				w,
				"No healthy backends available",
				http.StatusServiceUnavailable,
			)
			return
		}

		fmt.Println(
			"Proxying request to backend:",
			backend.URL.String(),
		)

		backend.ReverseProxy.ServeHTTP(w, r)
	})

	log.Printf("Load balancer running on port %d", config.Port)

	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%d", config.Port),
			nil,
		),
	)
}
