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
	ChangePassword(userID int64, currentPassword, newPassword string) error
	ResetPassword(email string) error
	DeleteToken(userID int64, tokenString string) error
	DeleteUser(id int64) error
}
