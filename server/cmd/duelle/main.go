package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rusilkoirala/duelle/internal/config"
	"github.com/rusilkoirala/duelle/internal/handler"
)

func main() {
	// get the Port for the server to run in
	cfg := config.Load()

	// the basic server
	mux := http.NewServeMux()

	mux.Handle("/", handler.NewStaticHandler(cfg.StaticDir))

	mux.Handle("/ws", handler.NewWSHandler())

	// health check point
	mux.HandleFunc("GET /health", handler.Health)

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("Server is running on http://localhost:%s", cfg.Port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server failed to start: %v", err)
	}
}
