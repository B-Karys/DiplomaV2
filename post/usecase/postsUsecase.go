package usecase

import (
	"DiplomaV2/post/models"
)

type PostUseCase interface {
	CreatePost(post *models.Post) error
	GetPostById(id int64) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	UpdatePost(id int64, authorID int64, title string, content string, postType string) error
	DeletePost(id int64) error
	GetFilteredPosts(authorID, postType string /*, additional filter criteria*/) ([]*models.Post, error)

	// Delete GetPostByAuthor(authorID int64) ([]*models.Post, error)
	// Delete GetPostByType(postType string) ([]*models.Post, error)
}
