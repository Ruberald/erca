package main

import (
	"log"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Backend struct {
	URL *url.URL
	Alive bool
	ReverseProxy *httputil.ReverseProxy
}

func newBackend(urlStr string) *Backend {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	return &Backend{
		URL: u,
		Alive: true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(u),
	}
}

func (b *Backend) IsAlive() bool {
	return b.Alive
}

func main() {
	targetURL, err := url.Parse("http://localhost:8081") // Backend server URL
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Proxying request to backend:", r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	log.Println("Load balancer is running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
