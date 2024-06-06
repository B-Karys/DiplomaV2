package repository

import (
	"DiplomaV2/backend/internal/database"
	"DiplomaV2/backend/internal/entity"
	"DiplomaV2/backend/post"
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

func (m *postRepository) Insert(post *entity.Post) error {
	result := m.DB.GetDb().Create(post)
	return result.Error
}

func (m *postRepository) GetByID(postID int64) (*entity.Post, error) {
	var thePost entity.Post
	if err := m.DB.GetDb().Where("id = ?", postID).First(&thePost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("thePost not found")
		}
		return nil, err
	}
	return &thePost, nil
}

func (m *postRepository) Delete(id int64) error {
	if id < 1 {
		return gorm.ErrRecordNotFound
	}

	result := m.DB.GetDb().Delete(&entity.Post{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (m *postRepository) Update(post *entity.Post) error {
	result := m.DB.GetDb().Save(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *postRepository) DeleteAllForUser(userID int64) error {
	if userID < 1 {
		return gorm.ErrRecordNotFound
	}

	result := m.DB.GetDb().Where("author_id = ?", userID).Delete(&entity.Post{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *postRepository) GetFilteredPosts(name, description, author, postType string, skills []string, filters post.Filters) ([]*entity.Post, post.Metadata, error) {
	var posts []*entity.Post
	query := m.DB.GetDb().Model(&entity.Post{})

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if description != "" {
		query = query.Where("description ILIKE ?", "%"+description+"%")
	}
	if author != "" {
		query = query.Where("author_id = ?", author)
	}
	if postType != "" {
		query = query.Where("type = ?", postType)
	}
	if len(skills) > 0 {
		query = query.Where("skills @> ?", pq.Array(skills))
	}

	var totalRecords int64
	countQuery := *query
	if err := countQuery.Count(&totalRecords).Error; err != nil {
		return nil, post.Metadata{}, err
	}

	if filters.Sort != "" && contains(filters.SortSafeList, filters.Sort) {
		query = query.Order(fmt.Sprintf("%s %s", filters.SortColumn(), filters.SortDirection()))
	}

	query = query.Offset((filters.Page - 1) * filters.PageSize).Limit(filters.PageSize)

	if err := query.Find(&posts).Error; err != nil {
		return nil, post.Metadata{}, err
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

func calculateMetadata(totalRecords, page, pageSize int) post.Metadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return post.Metadata{}
	}
	return post.Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
