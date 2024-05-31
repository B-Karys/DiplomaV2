package repository

import (
	"DiplomaV2/backend/internal/entity"
)

type UserRepository interface {
	Insert(user *entity.User) error
	GetByID(id int64) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	GetForToken(tokenScope, tokenPlaintext string) (*entity.User, error)
	Delete(id int64) error
}
