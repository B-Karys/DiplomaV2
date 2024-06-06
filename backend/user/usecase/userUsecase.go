package usecase

import (
	"DiplomaV2/backend/internal/entity"
	"mime/multipart"
)

type UserUseCase interface {
	Registration(user *entity.User) (*entity.Token, error)
	Activation(token string) error
	Authentication(user *entity.User) (string, error)
	GetUserById(id int64) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	UpdateUserInfo(int64, string, string, string, string, string, []string, string) error
	UploadProfileImage(userID int64, file *multipart.FileHeader) (string, error)
	ChangePassword(userID int64, currentPassword, newPassword string) error
	ForgotPassword(email string) (string, error)
	ResetPassword(string, string) error
	DeleteUser(id int64) error
}
