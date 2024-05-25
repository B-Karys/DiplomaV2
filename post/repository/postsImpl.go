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

func (m *postRepository) GetByAuthor(authorid int64) ([]*models.Post, error) {
	//if authorid < 1 {
	//	return nil, gorm.ErrRecordNotFound
	//}
	//
	var posts []*models.Post
	//if err := m.DB.GetDb().Where("authorid = ?", authorid).Find(&posts).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return nil, err
	//	}
	//	return nil, err
	//}

	return posts, nil
}

func (m *postRepository) DeleteAllForUser(userID int64) error {
	if userID < 1 {
		return gorm.ErrRecordNotFound
	}

	result := m.DB.GetDb().Where("authorid = ?", userID).Delete(&models.Post{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
