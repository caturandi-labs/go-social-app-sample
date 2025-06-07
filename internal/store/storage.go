package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	DatabaseQueryTimeout = 15 * time.Second
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("resource already exists")
)

type Storage struct {
	Posts interface {
		GetByID(ctx context.Context, id int64) (*Post, error)
		Create(context.Context, *Post) error
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(ctx context.Context, id int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(ctx context.Context, id int64) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerID int64, userID int64) error
		Unfollow(ctx context.Context, followerID int64, userID int64) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db: db},
		Users:     &UsersStore{db: db},
		Comments:  &CommentsStore{db: db},
		Followers: &FollowersStore{db: db},
	}
}
