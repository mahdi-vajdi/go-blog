package store

import (
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mahdi-vajdi/go-blog/internal/types"
	"log"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connectionString string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the postgres database")
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS posts (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL ,
	content TEXT,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)
	`

	_, err := s.db.Exec(createTableQuery)
	return err
}

func (s *PostgresStore) CreatePost(title, content string) (*types.Post, error) {
	query := `
	INSERT INTO posts (title, content) VALUES ($1, $2) RETURNING id, title, content, created_at
	`

	var post types.Post

	err := s.db.QueryRow(query, title, content).Scan(
		&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostgresStore) GetPosts() ([]types.Post, error) {
	query := `
	SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC
	`

	var posts []types.Post

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var post types.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return posts, nil
}

func (s *PostgresStore) GetPostByID(id int64) (*types.Post, error) {
	query := `
	SELECT id, title, content, created_at FROM posts WHERE id = $1
	`

	var post types.Post
	err := s.db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (s *PostgresStore) UpdatePost(id int64, title, content string) (*types.Post, error) {
	query := `
	UPDATE posts SET title = $1, content = $2 WHERE id = $3 RETURNING id, title, content, created_at
	`

	var post types.Post
	err := s.db.QueryRow(query, title, content, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (s *PostgresStore) DeletePost(id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}
