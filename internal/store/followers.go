package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `INSERT INTO followers(user_id, follower_id) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(ctx, DatabaseQueryTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("duplicate key value violates unique constraint %s", err)
		}
	}
	return nil
}

func (s *FollowersStore) Unfollow(ctx context.Context, followerID int64, userID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`
	ctx, cancel := context.WithTimeout(ctx, DatabaseQueryTimeout)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	return err

}
