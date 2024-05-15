package usecase

import (
	repository2 "DiplomaV2/token/repository"
	"DiplomaV2/user/models"
	"DiplomaV2/user/repository"
)

type userUseCaseImpl struct {
	repo repository.UserRepository
}

func (u *userUseCaseImpl) Activation(tokenPlaintext string) error {
	_, err := u.repo.GetForToken(repository2.ScopeActivation, tokenPlaintext)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUseCaseImpl) ResetPassword(email string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) Authentication(email, password string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) Registration(user *models.User) error {
	err := u.repo.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUseCaseImpl) UpdateUser(user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) DeleteUser(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) ShowOneUser(id int64) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCaseImpl{
		repo: repo,
	}
}
