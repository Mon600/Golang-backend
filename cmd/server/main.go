package main

import (
	"log"
	"net/http"
	"Golang-backend/internal/middleware"
)


func StartServer(srv *http.Server) {
	log.Println("Starting server on", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	handler := middleware.Logging(mux)

	srv := &http.Server{
		Addr: ":8000",
		Handler: handler,
	}
	go StartServer(srv)
	select {}
}
