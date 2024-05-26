package repository

import (
	"DiplomaV2/database"
	"DiplomaV2/post/models"
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type postRepository struct {
	DB database.Database
}

func NewPostRepository(db database.Database) PostRepository {
	return &postRepository{DB: db}
}

func (m *postRepository) Insert(post *models.Post) error {
	result := m.DB.GetDb().Create(post)
	return result.Error
}

func (m *postRepository) GetByID(postID int64) (*models.Post, error) {
	var post models.Post
	if err := m.DB.GetDb().Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (m *postRepository) Delete(id int64) error {
	if id < 1 {
		return gorm.ErrRecordNotFound
	}

	result := m.DB.GetDb().Delete(&models.Post{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (m *postRepository) Update(post *models.Post) error {
	result := m.DB.GetDb().Save(post)
	// Check for errors
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *postRepository) DeleteAllForUser(userID int64) error {
	if userID < 1 {
		return gorm.ErrRecordNotFound
	}

	result := m.DB.GetDb().Where("author_id = ?", userID).Delete(&models.Post{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *postRepository) GetFilteredPosts(name, description, author, postType string, skills []string, filters models.Filters) ([]*models.Post, error) {
	var posts []*models.Post
	query := m.DB.GetDb().Model(&models.Post{})

	// Apply filters
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

	// Apply sorting
	if filters.Sort != "" && contains(filters.SortSafeList, filters.Sort) {
		query = query.Order(fmt.Sprintf("%s %s", filters.SortColumn(), filters.SortDirection()))
	}

	// Apply pagination
	query = query.Offset((filters.Page - 1) * filters.PageSize).Limit(filters.PageSize)

	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
