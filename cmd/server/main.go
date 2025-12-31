package main

import (
	"Golang-backend/internal/data"
	"Golang-backend/internal/repository"
	"Golang-backend/internal/handlers"
	"Golang-backend/internal/middleware"
	"log"
	"net/http"

	_"github.com/golang-migrate/migrate/v4/database"
)

func StartServer(srv *http.Server) {
	log.Println("Starting server on", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed %v", err)
	}
}

func main() {
	db := data.NewDB("postgres://postgres:1@localhost:5432/go-backend?sslmode=disable")

	defer db.Close()
	
	repository.RunMigrations(db, "./internal/migrations")


	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)
	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.CreateUser)
	mux.HandleFunc("/user/", userHandler.GetUser)
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
