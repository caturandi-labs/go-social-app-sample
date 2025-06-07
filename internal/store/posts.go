package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type Post struct {
	ID        int64        `json:"id"`
	Content   string       `json:"content"`
	Title     string       `json:"title"`
	UserID    int64        `json:"user_id"`
	Version   int64        `json:"version"`
	Tags      []string     `json:"tags"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	Comments  []Comment    `json:"comments"`
	User      User         `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := "INSERT INTO posts (content, title, user_id, version, tags) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at;"

	ctx, cancel := context.WithTimeout(ctx, DatabaseQueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, post.Version, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, content, title, user_id, version,  tags, updated_at, created_at FROM posts WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, DatabaseQueryTimeout)
	defer cancel()

	var post Post

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		&post.Version,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return &post, nil
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := "UPDATE posts SET title = $1, content = $2, version = version + 1 WHERE id = $3 AND version = $4 RETURNING version;"

	ctx, cancel := context.WithTimeout(ctx, DatabaseQueryTimeout)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, post.Title, post.Content, post.ID, post.Version)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) Delete(ctx context.Context, postID int64) error {
	query := "DELETE FROM posts WHERE id = $1;"

	ctx, cancel := context.WithTimeout(ctx, DatabaseQueryTimeout)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err

	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) GetUserFeed(ctx context.Context, id int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
		SELECT
			p.id,p.user_id,p.title,p.content,p.created_at, p.version, p.tags, u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON p.user_id = f.follower_id
		WHERE f.user_id = $1 OR p.user_id = $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3;
	`
	ctx, cancel := context.WithTimeout(ctx, DatabaseQueryTimeout)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, id, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var feeds []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata
		err = rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, post)
	}
	return feeds, nil
}
