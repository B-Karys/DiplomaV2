package usecase

import (
	"DiplomaV2/user/models"
)

type UserUseCase interface {
	Registration(user *models.User) (*models.Token, error)
	Activation(token string) error
	Authentication(user *models.User) (string, error)
	GetUserById(id int64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	ChangePassword(user *models.User) error
	ResetPassword(email string) error
	DeleteUser(id int64) error
}
