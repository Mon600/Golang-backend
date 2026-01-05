package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"

	"Golang-backend/internal/model"
	"Golang-backend/internal/services"
)

type UserHandler struct {
	UserServ *services.UserService
}

func NewUserHandler(UserServ *services.UserService) *UserHandler {
	return &UserHandler{UserServ: UserServ}
}

func (h *UserHandler) CreateUser(r *http.Request, w http.ResponseWriter) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	user := model.User{Name: req.Name, Email: req.Email}
	id, err := h.UserServ.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   id,
		"name": req.Name,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "missing user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := h.UserServ.GetUser(id)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to get user: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPut {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var req model.UserUpdate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" && req.Email == "" {
		http.Error(w, "Nothing to change, fields is empty", http.StatusUnprocessableEntity)
		return
	}

	defer r.Body.Close()

	new_data, err := h.UserServ.UpdateUser(id, req)
	if err != nil {
		http.Error(w, string(err.Error()), http.StatusInternalServerError)
		return
	}
	if new_data == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(new_data)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	deleted_email, err := h.UserServ.DeleteUser(id)
	if err != nil {
		log.Fatal("error: %w", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	if deleted_email == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deleted_email)
}
