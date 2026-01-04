package main

import (
	"Golang-backend/internal/data"
	"Golang-backend/internal/handlers"
	"Golang-backend/internal/middleware"
	"Golang-backend/internal/repository"
	"Golang-backend/internal/services"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/database"
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
	userServ := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userServ)
	r := chi.NewRouter()
	r.Use(middleware.Logging)
	r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		userHandler.CreateUser(r, w)
	})
	r.Put("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		userHandler.UpdateUser(w, r)
	})
	r.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		userHandler.GetUser(w, r)
	})
	r.Get("/main", func(w http.ResponseWriter, r *http.Request) {
		handlers.MainHandler(w, r)
	})

	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	go StartServer(srv)
	select {}
}
