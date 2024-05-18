package usecase

import (
	tokenRepository "DiplomaV2/token/repository"
	userRepository "DiplomaV2/user/repository"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"os"
	"time"
)

type tokenUseCaseImpl struct {
	tokenRepo tokenRepository.TokenRepository
	userRepo  userRepository.UserRepository
}

var (
	InvalidPassword = errors.New("Invalid password")
)

func (t tokenUseCaseImpl) createActivationToken(email string) error {
	//token, err := t.tokenRepo.New(user.ID, 1*time.Hour, repository.ScopeActivation)
	//if token == nil || err != nil {
	//	return "Token is not created", err
	//}
	//return token.Plaintext, err
	//TODO implement me
	panic("implement me")
}

func (t tokenUseCaseImpl) createAuthenticationToken(email, password string) (string, error) {
	user, err := t.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "User email is not found", err
		}
		return "Server error", err
	}

	match, err := user.Password.Matches(password)
	if err != nil {
		return "Error is: ", err
	}
	if !match {
		return "Password didn't match", InvalidPassword
	}

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

	// Optionally create a record in the token repository
	_, err = t.tokenRepo.New(user.ID, 1*time.Hour, tokenRepository.ScopeAuthentication)
	if err != nil {
		return "Error while Inserting token to DB", err
	}

	return jwtToken, nil
}

func (t tokenUseCaseImpl) createPasswordResetToken(email string) error {
	//TODO implement me
	panic("implement me")
}

func NewTokenUseCase(tokenRepo tokenRepository.TokenRepository, userRepo userRepository.UserRepository) TokenUseCase {
	return &tokenUseCaseImpl{
		tokenRepo: tokenRepo,
		userRepo:  userRepo,
	}
}
