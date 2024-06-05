package usecase

import (
	"DiplomaV2/backend/internal/entity"
	"DiplomaV2/backend/post"
)

type PostUseCase interface {
	CreatePost(post *entity.Post) error
	GetPostById(id int64) (*entity.Post, error)
	UpdatePost(id int64, authorID int64, title string, content string, skills []string, postType string) error
	DeletePost(id int64) error
	GetFilteredPosts(name, description, authorID, postType string, skills []string, filters post.Filters) ([]*entity.Post, post.Metadata, error)
}
