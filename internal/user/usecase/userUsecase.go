package usecase

import "DiplomaV2/domain/models"

type UserUseCase interface {
	Registration(input *struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}) error
	Activation(token string) error
	ForgotPassword(email string) error
	Authenticate(email, password string) error
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
}
