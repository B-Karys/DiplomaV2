package repository

import (
	"DiplomaV2/database"
	"DiplomaV2/post/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type postRepository struct {
	DB database.Database
}

func NewPostRepository(db database.Database) PostRepository {
	return &postRepository{DB: db}
}

/*
func (m *postRepository) GetByType(postType string) ([]*models.Post, error) {
	var posts []*models.Post
	// Assuming your Post model is named Post and your table name is "posts"
	if err := m.DB.GetDb().Where("type = ?", postType).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
*/

func (m *postRepository) Insert(post *models.Post) error {
	result := m.DB.GetDb().Create(post)
	return result.Error
}

func (m *postRepository) GetAll() ([]*models.Post, error) {
	var posts []*models.Post
	if err := m.DB.GetDb().Find(&posts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return posts, nil
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

/*
func (m *postRepository) GetByAuthor(authorid int64) ([]*models.Post, error) {
	if authorid < 1 {
		return nil, gorm.ErrRecordNotFound
	}

	var posts []*models.Post
	if err := m.DB.GetDb().Where("author_id = ?", authorid).Find(&posts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return posts, nil
}
*/

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

func (m *postRepository) GetByAuthorAndType(authorID string, postType string) ([]*models.Post, error) {
	var posts []*models.Post
	// Assuming your Post model is named Post and your table name is "posts"
	if authorID != "" {
		query := m.DB.GetDb().Where("author_id = ?", authorID)
		if postType != "" {
			// If postType is provided, add it as a filter
			query = query.Where("type = ?", postType)
		}
		if err := query.Find(&posts).Error; err != nil {
			return nil, err
		}
	} else if postType != "" {
		// If only postType is provided, filter by postType
		if err := m.DB.GetDb().Where("type = ?", postType).Find(&posts).Error; err != nil {
			return nil, err
		}
	} else {
		// If neither authorID nor postType is provided, return all posts
		if err := m.DB.GetDb().Find(&posts).Error; err != nil {
			return nil, err
		}
	}
	return posts, nil
}
