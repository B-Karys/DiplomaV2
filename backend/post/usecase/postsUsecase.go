package usecase

import (
	"DiplomaV2/backend/internal/entity"
	postsFilter "DiplomaV2/backend/post"
)

type PostUseCase interface {
	CreatePost(post *entity.Post) error
	GetPostById(id int64) (*entity.Post, error)
	UpdatePost(id int64, authorID int64, post *entity.Post) error
	DeletePost(id int64) error
	GetFilteredPosts(post *entity.Post, filters postsFilter.Filters) ([]*entity.Post, postsFilter.Metadata, error)
}
