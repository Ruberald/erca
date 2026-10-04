package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type Config struct {
	Backends []string `json:"backends"`
}

func startBackend(backendURL string) {
	u, err := url.Parse(backendURL)
	if err != nil {
		log.Printf("Invalid backend URL %q: %v", backendURL, err)
		return
	}

	port := u.Port()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from backend %s\n", port)
	})

	log.Printf("Starting backend on %s", backendURL)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Printf("Backend %s stopped: %v", port, err)
	}
}

func main() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, backendURL := range config.Backends {
		wg.Add(1)

		go func(backendURL string) {
			defer wg.Done()
			startBackend(backendURL)
		}(backendURL)
	}

	wg.Wait()
}
