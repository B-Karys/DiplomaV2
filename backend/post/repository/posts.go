package repository

import (
	"DiplomaV2/backend/post/models"
)

type PostRepository interface {
	Insert(post *models.Post) error
	GetByID(id int64) (*models.Post, error)
	Delete(id int64) error
	Update(post *models.Post) error
	DeleteAllForUser(authorid int64) error
	GetFilteredPosts(name, description, author, postType string, skills []string, filters models.Filters) ([]*models.Post, error)

	// Delete GetByAuthor(authorid int64) ([]*models.Post, error)
	// Delete GetByType(postType string) ([]*models.Post, error)
}
