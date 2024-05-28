package usecase

import (
	"DiplomaV2/backend/post/models"
)

type PostUseCase interface {
	CreatePost(post *models.Post) error
	GetPostById(id int64) (*models.Post, error)
	UpdatePost(id int64, authorID int64, title string, content string, skills []string, postType string) error
	DeletePost(id int64) error
	GetFilteredPosts(name, description, author, postType string, skills []string, filters models.Filters) ([]*models.Post, error)

	// Delete GetPostByAuthor(authorID int64) ([]*models.Post, error)
	// Delete GetPostByType(postType string) ([]*models.Post, error)
}
