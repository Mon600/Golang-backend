package repository

import (
	"database/sql"
	"fmt"
	"Golang-backend/internal/model"
)


type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user model.User) (int, error) {
	const query = `
	INSERT INTO users (name, email)
	VALUES ($1, $2)
	RETURNING id
	`
	var id int
	err := r.db.QueryRow(query, user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	const query = `
	SELECT id, name, email, created_at
	FROM users
	WHERE users.id = $1
	`

	row := r.db.QueryRow(query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user %w", err)
	}
	return &user, nil
}