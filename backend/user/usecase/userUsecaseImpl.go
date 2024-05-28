package usecase

import (
	"DiplomaV2/backend/internal/validator"
	"DiplomaV2/backend/user/models"
	"DiplomaV2/backend/user/repository"
	"DiplomaV2/backend/user/tokenRepository"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"os"
	"time"
)

type userUseCaseImpl struct {
	repo      repository.UserRepository
	tokenRepo tokenRepository.TokenRepository
}

var (
	TokenCreationFailed = errors.New("Token creation failed")
	ErrWrongPassword    = errors.New("Wrong password")
	InvalidToken        = errors.New("Token is invalid")
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

func (u *userUseCaseImpl) ChangePassword(userId int64, currentPassword string, newPassword string) error {
	user, err := u.repo.GetByID(userId)
	if err != nil {
		return err
	}

	match, err := user.Password.Matches(currentPassword)
	if err != nil {
		return err
	}

	if !match {
		return ErrWrongPassword
	}

	v := validator.New()
	validator.ValidatePasswordPlaintext(v, newPassword)
	if !v.Valid() {
		return errors.New("password should contain more than 8 characters")
	}

	err = user.Password.Set(newPassword)
	if err != nil {
		return err
	}
	err = u.repo.Update(user)
	if err != nil {
		return err
	}

	return err
}

func (u *userUseCaseImpl) Authentication(user *models.User) (string, error) {
	token, err := u.createAuthenticationToken(user)
	if err != nil {
		return TokenCreationFailed.Error(), err
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

func (u *userUseCaseImpl) UpdateUserInfo(userID int64, name, surname string, telegram string, discord string, skills []string, profileImage string) error {
	user, err := u.repo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Name = name
	user.Surname = surname
	user.Telegram = telegram
	user.Discord = discord
	user.Skills = skills
	user.ProfileImage = profileImage

	err = u.repo.Update(user)
	if err != nil {
		return err
	}
	return err
}

func (u *userUseCaseImpl) DeleteUser(id int64) error {
	err := u.repo.Delete(id)
	if err != nil {
		return err
	}
	return err
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
func (u *userUseCaseImpl) ForgotPassword(email string) (string, error) {
	// Find the user by email
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	// Generate the password reset token
	token, err := u.tokenRepo.New(user.ID, 24*time.Hour, tokenRepository.ScopePasswordReset)
	if err != nil {
		return "", err
	}

	return token.Plaintext, nil
}

func (u *userUseCaseImpl) ResetPassword(tokenString, newPassword string) error {

	v := validator.New()
	user, err := u.repo.GetForToken(tokenRepository.ScopePasswordReset, tokenString)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			v.AddError("token", "invalid or expired password reset token")
			return InvalidToken
		default:
			return err
		}
	}

	// Update the user's password
	err = user.Password.Set(newPassword)
	if err != nil {
		return err
	} // Assume hashing is done in repository

	// Save the updated user
	err = u.repo.Update(user)
	if err != nil {
		return err
	}

	// Delete the token after successful password reset
	err = u.tokenRepo.DeleteAllForUser(tokenRepository.ScopePasswordReset, user.ID)
	if err != nil {
		return err
	}

	return nil
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
		return TokenCreationFailed.Error(), err
	}

	// Optionally create a record in the token tokenRepository
	//_, err = u.tokenRepo.New(user.ID, 1*time.Hour, tokenRepository.ScopeAuthentication)
	//if err != nil {
	//	return "Error while Inserting token to DB", err
	//}

	return jwtToken, nil
}

func NewUserUseCase(repo repository.UserRepository, tokenRepo tokenRepository.TokenRepository) UserUseCase {
	return &userUseCaseImpl{
		repo:      repo,
		tokenRepo: tokenRepo,
	}
}
