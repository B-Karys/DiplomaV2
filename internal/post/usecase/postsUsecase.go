package usecase

import "DiplomaV2/domain/models"

type PostUseCase interface {
	CreatePost(post *models.Post) error
	ShowOncePost(id int64) (*models.Post, error)
	DeletePost(id int64) error
	UpdatePost(post *models.Post) error
}
