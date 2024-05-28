package usecase

import (
	"DiplomaV2/backend/user/models"
)

type UserUseCase interface {
	Registration(user *models.User) (*models.Token, error)
	Activation(token string) error
	Authentication(user *models.User) (string, error)
	GetUserById(id int64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUserInfo(int64, string, string, string, string, []string, string) error
	ChangePassword(userID int64, currentPassword, newPassword string) error
	ForgotPassword(email string) (string, error)
	ResetPassword(string, string) error
	DeleteUser(id int64) error
}
