package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from backend at port %s!", port)
	})

	fmt.Println("Backend server is running on port", port)

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
