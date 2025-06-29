package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64        `json:"id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Password  string       `json:"-"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {

	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at"

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(
		&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}
	return nil

}

func (s *UsersStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := "SELECT id, username, email, created_at FROM users WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, DatabaseQueryTimeout)
	defer cancel()

	u := &User{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return u, nil

}
