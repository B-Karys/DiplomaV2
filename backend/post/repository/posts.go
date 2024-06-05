package repository

import (
	"DiplomaV2/backend/internal/entity"
	"DiplomaV2/backend/post"
)

type PostRepository interface {
	Insert(post *entity.Post) error
	GetByID(id int64) (*entity.Post, error)
	Delete(id int64) error
	Update(post *entity.Post) error
	DeleteAllForUser(authorid int64) error
	GetFilteredPosts(name, description, author, postType string, skills []string, filters post.Filters) ([]*entity.Post, post.Metadata, error)
}
