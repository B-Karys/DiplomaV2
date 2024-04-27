package usecase

import "DiplomaV2/domain/models"

type UserUseCase interface {
	Registration(user *models.User) error
	Activation(token string) error
	ForgotPassword(email string) error
	Authenticate(email, password string) error
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
}
