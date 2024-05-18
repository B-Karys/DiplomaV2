package usecase

import (
	"DiplomaV2/user/models"
	"DiplomaV2/user/repository"
	"DiplomaV2/user/tokenRepository"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"os"
	"time"
)

type userUseCaseImpl struct {
	repo      repository.UserRepository
	tokenRepo tokenRepository.TokenRepository
}

var (
	InvalidPassword     = errors.New("Invalid password")
	TokenCreationFailed = errors.New("Token creation failed")
)

func (u *userUseCaseImpl) Activation(tokenPlaintext string) error {
	user, err := u.repo.GetForToken(tokenRepository.ScopeActivation, tokenPlaintext)
	if err != nil {
		return err
	}
	user.Activated = true
	err = u.repo.Update(user)
	if err != nil {
		return err
	}

	err = u.tokenRepo.DeleteAllForUser(tokenRepository.ScopeActivation, user.ID)
	if err != nil {
		return err
	}
	return err
}

func (u *userUseCaseImpl) ChangePassword(user *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) ResetPassword(email string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userUseCaseImpl) Authentication(user *models.User) (string, error) {
	token, err := u.createAuthenticationToken(user)
	if err != nil {
		return "Token is not created", err
	}
	return token, nil
}

func (u *userUseCaseImpl) Registration(user *models.User) (*models.Token, error) {
	err := u.repo.Insert(user)
	if err != nil {
		return nil, err
	}

	token, err := u.createActivationToken(user)
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

func (u *userUseCaseImpl) GetUserById(id int64) (*models.User, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCaseImpl) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCaseImpl) createActivationToken(user *models.User) (*models.Token, error) {
	token, err := u.tokenRepo.New(user.ID, 1*time.Hour, tokenRepository.ScopeActivation)
	if token == nil || err != nil {
		return nil, err
	}
	return token, nil
}

func (u *userUseCaseImpl) createAuthenticationToken(user *models.User) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iss": "TeamFinder",
		"aud": "TeamFinder",
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "Error while Creating Token", err
	}

	// Optionally create a record in the token tokenRepository
	_, err = u.tokenRepo.New(user.ID, 1*time.Hour, tokenRepository.ScopeAuthentication)
	if err != nil {
		return "Error while Inserting token to DB", err
	}

	return jwtToken, nil
}

func (u *userUseCaseImpl) createPasswordResetToken(email string) error {
	//TODO implement me
	panic("implement me")
}

func NewUserUseCase(repo repository.UserRepository, tokenRepo tokenRepository.TokenRepository) UserUseCase {
	return &userUseCaseImpl{
		repo:      repo,
		tokenRepo: tokenRepo,
	}
}
