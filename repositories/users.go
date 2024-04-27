package repositories

import "DiplomaV2/domain/models"

type UsersRepository interface {
	Insert(user *models.User) error
	GetByID(id int64) (*models.User, error)
	Get(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	GetForToken(tokenScope, tokenPlaintext string) (*models.User, error)
	Delete(id int64) error
}
