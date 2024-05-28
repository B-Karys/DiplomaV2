package repository

import (
	"DiplomaV2/backend/user/models"
)

type UserRepository interface {
	Insert(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	GetForToken(tokenScope, tokenPlaintext string) (*models.User, error)
	Delete(id int64) error
}
