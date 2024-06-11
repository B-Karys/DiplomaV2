package usecase

import (
	"DiplomaV2/backend/internal/entity"
	"DiplomaV2/backend/internal/helpers"
	"DiplomaV2/backend/internal/validator"
	"DiplomaV2/backend/user/repository"
	"DiplomaV2/backend/user/tokenRepository"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

type userUseCaseImpl struct {
	repo      repository.UserRepository
	tokenRepo tokenRepository.TokenRepository
}

func (u *userUseCaseImpl) GetAllUsers() ([]*entity.User, error) {
	users, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
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

func (u *userUseCaseImpl) Authentication(user *entity.User) (string, error) {
	token, err := u.createAuthenticationToken(user)
	if err != nil {
		return TokenCreationFailed.Error(), err
	}
	return token, nil
}

func (u *userUseCaseImpl) Registration(user *entity.User) (*entity.Token, error) {
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

func (u *userUseCaseImpl) UpdateUserInfo(userID int64, name, surname, username, telegram, discord string, skills []string, profileImage string) error {
	user, err := u.repo.GetByID(userID)
	if err != nil {
		return err
	}

	user.Name = name
	user.Surname = surname
	user.Username = username
	user.Telegram = telegram
	user.Discord = discord
	user.Skills = skills
	user.Version = user.Version + 1

	if profileImage != "" {
		user.ProfileImage = profileImage
	}

	err = u.repo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUseCaseImpl) UploadProfileImage(userID int64, file *multipart.FileHeader) (string, error) {
	user, err := u.repo.GetByID(userID)
	if err != nil {
		return "", err
	}

	fmt.Println("Uploading profile image...")

	src, err := file.Open()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", errors.New("Failed to open file")
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			return
		}
	}(src)

	// Initialize GCS client
	ctx := context.Background()
	client, err := helpers.NewStorageClient(ctx, "C:/Users/krump/Downloads/lucid-volt-424719-f0-5df86076a210.json")
	if err != nil {
		return "", errors.New("Failed to initialize GCS client")
	}

	objectName := fmt.Sprintf("%d/%d", userID, user.Version) // Unique object name based on user ID
	if err := helpers.UploadFileToGCS(ctx, client, "teamfinderimages", objectName, src); err != nil {
		return "", errors.New("Failed to upload file to GCS")
	}

	flood := user.ProfileImage == "https://storage.googleapis.com/teamfinderimages/default_photo.png"
	baseUrl := "https://storage.googleapis.com/teamfinderimages/" + strconv.FormatInt(user.ID, 10) + "/"
	var profileImageURL string

	if flood {
		// If the user's profile image is default, set profile image URL to baseUrl/1
		profileImageURL = baseUrl + "1"
	} else {
		// If the user's profile image is not default, extract the number from the URL, increment it by 1, and construct the new URL
		endIndex := strings.LastIndex(user.ProfileImage, "/") + 1
		numStr := user.ProfileImage[endIndex:]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			// Handle error
			return "", err
		}
		num++
		profileImageURL = baseUrl + strconv.Itoa(num)
	}

	return profileImageURL, nil
}

func (u *userUseCaseImpl) DeleteUser(id int64) error {
	err := u.repo.Delete(id)
	if err != nil {
		return err
	}
	return err
}

func (u *userUseCaseImpl) GetUserById(id int64) (*entity.User, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCaseImpl) GetUserByEmail(email string) (*entity.User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *userUseCaseImpl) ForgotPassword(email string) (string, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

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

	err = user.Password.Set(newPassword)
	if err != nil {
		return err
	}

	err = u.repo.Update(user)
	if err != nil {
		return err
	}

	err = u.tokenRepo.DeleteAllForUser(tokenRepository.ScopePasswordReset, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCaseImpl) createActivationToken(user *entity.User) (*entity.Token, error) {
	token, err := u.tokenRepo.New(user.ID, 1*time.Hour, tokenRepository.ScopeActivation)
	if token == nil || err != nil {
		return nil, err
	}
	return token, nil
}

func (u *userUseCaseImpl) createAuthenticationToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iss": "TeamFinder",
		"aud": "TeamFinder",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return TokenCreationFailed.Error(), err
	}

	return jwtToken, nil
}

func NewUserUseCase(repo repository.UserRepository, tokenRepo tokenRepository.TokenRepository) UserUseCase {
	return &userUseCaseImpl{
		repo:      repo,
		tokenRepo: tokenRepo,
	}
}
