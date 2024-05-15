package usecase

import (
	"DiplomaV2/post/models"
)

type PostUseCase interface {
	CreatePost(post *models.Post) error
	ShowOnePost(id int64) (*models.Post, error)
	DeletePost(id int64) error
	UpdatePost(post *models.Post) error
}
