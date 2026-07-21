package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rusilkoirala/duelle/internal/handler"
)

func main() {
	// get the Port for the server to run in
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// the basic server
	mux := http.NewServeMux()

	// i made frontend to load from backend / route like give static files
	staticDir := http.Dir("./static")
	staticHandler := http.FileServer(staticDir)
	mux.Handle("/", staticHandler)

	// health check point
	mux.HandleFunc("GET /health", handler.Health)

	addr := fmt.Sprintf(":%s", port)

	log.Printf("Server is running on http://localhost:%s", port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server failed to start: %v", err)
	}
}
