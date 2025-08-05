package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "3001", "Port to run the WebSocket server on")
	flag.Parse()
	// Initialize the hub and start it in a goroutine
	hub := NewHub()
	go hub.Run()

	ws := &Server{Hub: hub}
	server := http.Server{
		Addr:    ":" + *port,
		Handler: ws,
	}

	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
