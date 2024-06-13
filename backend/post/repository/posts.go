package repository

import (
	"DiplomaV2/backend/internal/entity"
	postsFilter "DiplomaV2/backend/post"
)

type PostRepository interface {
	Insert(post *entity.Post) error
	GetByID(id int64) (*entity.Post, error)
	Delete(id int64) error
	Update(post *entity.Post) error
	DeleteAllForUser(authorID int64) error
	GetFilteredPosts(post *entity.Post, filters postsFilter.Filters) ([]*entity.Post, postsFilter.Metadata, error)
}
