package data

import (
	"DiplomaV2/domain/models"
	"DiplomaV2/internal/post/repository"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type postRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) repository.PostRepository {
	return &postRepository{DB: db}
}

func (m postRepository) Insert(post *models.Post) error {
	query := `
		INSERT INTO posts(name, description, authorid, type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	return m.DB.QueryRow(query, &post.Name, &post.Description, &post.AuthorID, &post.Type).Scan(&post.ID, &post.CreatedAt)
}

func (m postRepository) GetAll() ([]*models.Post, error) {
	query := `
		SELECT *
		FROM posts`

	rows, err := m.DB.Query(query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, gorm.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.CreatedAt,
			&post.Name,
			&post.Description,
			&post.AuthorID,
			&post.Type,
			&post.Version,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// method for fetching a specific record from the movies table.
func (m postRepository) GetByID(id int64) (*models.Post, error) {
	if id < 1 {
		return nil, gorm.ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM posts
		WHERE id = $1`

	var post models.Post

	err := m.DB.QueryRow(query, id).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Name,
		&post.Description,
		&post.AuthorID,
		&post.Type,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, gorm.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (m postRepository) Delete(id int64) error {
	if id < 1 {
		return gorm.ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM posts
		WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (m postRepository) Update(post *models.Post) error {
	// Add the 'AND version = $6' clause to the SQL query.
	query := `
UPDATE posts
SET name = $1, description = $2, type = $3, authorid = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	args := []any{
		post.Name,
		post.Description,
		post.Type,
		post.AuthorID,
		post.ID,
		post.Version, // Add the expected post version.
	}
	// Execute the SQL query. If no matching row could be found, we know the post
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	err := m.DB.QueryRow(query, args...).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return gorm.ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (m postRepository) GetByAuthor(authorid int64) ([]*models.Post, error) {
	if authorid < 1 {
		return nil, gorm.ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM posts
		WHERE authorid = $1`

	rows, err := m.DB.Query(query, authorid)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, gorm.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.CreatedAt,
			&post.Name,
			&post.Description,
			&post.AuthorID,
			&post.Type,
			&post.Version,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m postRepository) DeleteAllForUser(userID int64) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, userID)
	return err
}
