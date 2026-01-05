package repository

import (
	"Golang-backend/internal/model"
	"database/sql"
	"fmt"
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

func (r *UserRepository) GetByID(id int64) (*model.User, error) {
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

func (r *UserRepository) UpdateByID(id int64, data model.User) (*model.User, error) {
	const query = `
	UPDATE users
	SET name = $1,
		email = $2
	WHERE users.id = $3
	RETURNING id, name, email, created_at
	`
	row := r.db.QueryRow(query, data.Name, data.Email, id)
	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update user %w", err)
	}
	return &user, nil
}

func (r *UserRepository) DeleteByID(id int64) (*string, error) {
	const query = `
	DELETE FROM users
	WHERE users.id = $1
	RETURNING email`

	row := r.db.QueryRow(query, id)
	var email string
	err := row.Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to delet user with id %d, error: %w", id, err)
	}
	return &email, nil
}
