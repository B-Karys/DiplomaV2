package usecase

import (
	"DiplomaV2/user/models"
)

type UserUseCase interface {
	Registration(user *models.User) (*models.Token, error)
	Activation(token string) error
	Authentication(email, password string) error
	ResetPassword(email string) error
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
	ShowUserById(id int64) (*models.User, error)
	ShowUserByEmail(email string) (*models.User, error)
}
