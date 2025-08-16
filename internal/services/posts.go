package services

import (
	"database/sql"
	"time"

	"github.com/scopophobic/scopoBlog/internal/models"
)

func CreatePost(db *sql.DB, post *models.Post) (*models.Post, error) {
	post.Status = "draft"
	post.Visible = false
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	query := `INSERT INTO posts (title, slug, content, status, visible, created_at, updated_at)
			   VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := db.Exec(query, post.Title, post.Slug, post.Content, post.Status, post.Visible, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	post.ID = int(id)

	return post, nil

}

func GetAllPublishedPosts(db *sql.DB) ([]models.Post, error) {
	query := `SELECT id, title, slug, content, created_at, updated_at FROM posts
			   WHERE status = 'published' AND visible = TRUE ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		// Note: We are not scanning status or visible since we filtered by them.
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetPostBySlug retrieves a single published and visible post by its slug.
func GetPostBySlug(db *sql.DB, slug string) (*models.Post, error) {
	query := `SELECT id, title, slug, content, created_at, updated_at FROM posts
			   WHERE slug = ? AND status = 'published' AND visible = TRUE`

	row := db.QueryRow(query, slug)

	var post models.Post
	if err := row.Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil to indicate not found, but not an error.
		}
		return nil, err
	}
	return &post, nil
}

// i do not understand honestly
func UpdatePost(db *sql.DB, id int, post *models.Post) (*models.Post, error) {
	post.UpdatedAt = time.Now()
	query := `UPDATE posts SET title = ?, content = ?, status = ?, visible = ?, updated_at = ?
			   WHERE id = ?`
	_, err := db.Exec(query, post.Title, post.Content, post.Status, post.Visible, post.UpdatedAt, id)
	if err != nil {
		return nil, err
	}
	post.ID = id
	return post, nil
}

func GetAllDraftPost(db *sql.DB) ([]models.Post, error) {
	var query = `SELECT id, title, slug, content, created_at, updated_at FROM posts WHERE status = 'draft' AND visible = false ORDER BY created_at DESC`

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		// Note: We are not scanning status or visible since we filtered by them.
		if err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func DeletePost(db *sql.DB, id int) error {
	query := `DELETE FROM posts WHERE id = ?`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
