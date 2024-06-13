package repository

import (
	"DiplomaV2/backend/internal/database"
	"DiplomaV2/backend/internal/entity"
	postsFilter "DiplomaV2/backend/post"
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"math"
)

type postRepository struct {
	DB database.Database
}

func NewPostRepository(db database.Database) PostRepository {
	return &postRepository{DB: db}
}

func (r *postRepository) Insert(post *entity.Post) error {
	result := r.DB.GetDb().Create(post)
	return result.Error
}

func (r *postRepository) GetByID(postID int64) (*entity.Post, error) {
	var post entity.Post
	if err := r.DB.GetDb().Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Delete(id int64) error {
	if id < 1 {
		return gorm.ErrRecordNotFound
	}

	result := r.DB.GetDb().Delete(&entity.Post{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *postRepository) Update(post *entity.Post) error {
	result := r.DB.GetDb().Save(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *postRepository) DeleteAllForUser(userID int64) error {
	if userID < 1 {
		return gorm.ErrRecordNotFound
	}

	result := r.DB.GetDb().Where("author_id = ?", userID).Delete(&entity.Post{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *postRepository) GetFilteredPosts(post *entity.Post, filters postsFilter.Filters) ([]*entity.Post, postsFilter.Metadata, error) {
	var posts []*entity.Post
	query := r.DB.GetDb().Model(&entity.Post{})

	if post.Name != "" {
		query = query.Where("name ILIKE ?", "%"+post.Name+"%")
	}
	if post.Description != "" {
		query = query.Where("description ILIKE ?", "%"+post.Description+"%")
	}
	if post.AuthorID != 0 {
		query = query.Where("author_id = ?", post.AuthorID)
	}
	if post.Type != "" {
		query = query.Where("type = ?", post.Type)
	}
	if len(post.Skills) > 0 {
		query = query.Where("skills @> ?", pq.Array(post.Skills))
	}

	var totalRecords int64
	countQuery := *query
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, postsFilter.Metadata{}, err
	}

	if filters.Sort != "" && contains(filters.SortSafeList, filters.Sort) {
		query = query.Order(fmt.Sprintf("%s %s", filters.SortColumn(), filters.SortDirection()))
	}

	query = query.Offset((filters.Page - 1) * filters.PageSize).Limit(filters.PageSize)

	if err := query.Find(&posts).Error; err != nil {
		return nil, postsFilter.Metadata{}, err
	}

	metadata := calculateMetadata(int(totalRecords), filters.Page, filters.PageSize)
	return posts, metadata, nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func calculateMetadata(totalRecords, page, pageSize int) postsFilter.Metadata {
	if totalRecords == 0 {
		return postsFilter.Metadata{}
	}
	return postsFilter.Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
