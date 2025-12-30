package main

import (
	"Golang-backend/internal/handlers"
	"Golang-backend/internal/middleware"
	"log"
	"net/http"
)

func StartServer(srv *http.Server) {
	log.Println("Starting server on", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/main", handlers.MainHandler)
	mux.HandleFunc("/post", handlers.PostHandler)

	handler := middleware.Logging(mux)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: handler,
	}
	go StartServer(srv)
	select {}
}
