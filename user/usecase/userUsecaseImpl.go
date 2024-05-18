package usecase

import (
	repository2 "DiplomaV2/token/repository"
	"DiplomaV2/user/models"
	"DiplomaV2/user/repository"
	"github.com/pkg/errors"
	"time"
)

type userUseCaseImpl struct {
	repo      repository.UserRepository
	tokenRepo repository2.TokenRepository
}

var (
	InvalidPassword = errors.New("Invalid password")
)

func (u *userUseCaseImpl) Activation(tokenPlaintext string) error {
	user, err := u.repo.GetForToken(repository2.ScopeActivation, tokenPlaintext)
	if err != nil {
		return err
	}
	user.Activated = true
	err = u.repo.Update(user)
	if err != nil {
		return err
	}

	err = u.tokenRepo.DeleteAllForUser(repository2.ScopeActivation, user.ID)
	if err != nil {
		return err
	}
	return err
}

func (u *userUseCaseImpl) ResetPassword(email string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) Authentication(email, password string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) Registration(user *models.User) (*models.Token, error) {
	token, err := u.createActivationToken(user)
	if err != nil {
		return nil, err
	}

	err = u.repo.Insert(user)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (u *userUseCaseImpl) UpdateUser(user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) DeleteUser(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) ShowUserById(id int64) (*models.User, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCaseImpl) ShowUserByEmail(email string) (*models.User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCaseImpl) createActivationToken(user *models.User) (*models.Token, error) {
	token, err := u.tokenRepo.New(user.ID, 1*time.Hour, repository2.ScopeActivation)
	if token == nil || err != nil {
		return nil, err
	}
	return token, nil
}

//func (u *userUseCaseImpl) createAuthenticationToken(user *models.User) (string, error) {
//	match, err := user.Password.Matches(password)
//	if err != nil {
//		return "Error is: ", err
//	}
//	if !match {
//		return "Password didn't match", InvalidPassword
//	}
//
//	// Define JWT claims
//	claims := jwt.MapClaims{
//		"sub": user.ID,
//		"iat": time.Now().Unix(),
//		"nbf": time.Now().Unix(),
//		"exp": time.Now().Add(24 * time.Hour).Unix(),
//		"iss": "TeamFinder",
//		"aud": "TeamFinder",
//	}
//
//	// Create and sign the token
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
//	jwtToken, err := token.SignedString(jwtSecret)
//	if err != nil {
//		return "Error while Creating Token", err
//	}
//
//	// Optionally create a record in the token repository
//	_, err = t.tokenRepo.New(user.ID, 1*time.Hour, tokenRepository.ScopeAuthentication)
//	if err != nil {
//		return "Error while Inserting token to DB", err
//	}
//
//	return jwtToken, nil
//}

func (u *userUseCaseImpl) createPasswordResetToken(email string) error {
	//TODO implement me
	panic("implement me")
}

func NewUserUseCase(repo repository.UserRepository, tokenRepo repository2.TokenRepository) UserUseCase {
	return &userUseCaseImpl{
		repo:      repo,
		tokenRepo: tokenRepo,
	}
}
