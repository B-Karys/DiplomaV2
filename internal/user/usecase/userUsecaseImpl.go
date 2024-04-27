package usecase

//
//import (
//	"DiplomaV2/domain/models"
//	"DiplomaV2/internal/user/repository"
//	"errors"
//	"fmt"
//)
//
//type userUseCaseImpl struct {
//	repo repository.UserRepository
//}
//
//func NewUserUseCase(repo repository.UserRepository) UserUseCase {
//	return &userUseCaseImpl{
//		repo: repo,
//	}
//}
//
//func (u *userUseCaseImpl) Registration(user *models.User) error {
//	// Validate user data, e.g., email format, password strength, etc.
//	if !isValidEmail(user.Email) {
//		return errors.New("invalid email address")
//	}
//	// Additional validation checks...
//
//	// Check if the email is already registered
//	existingUser, err := u.repo.GetByEmail(user.Email)
//	if err == nil && existingUser != nil {
//		return errors.New("email already registered")
//	}
//
//	// Hash user password before saving to the database
//	hashedPassword, err := hashPassword(user.Password)
//	if err != nil {
//		return fmt.Errorf("failed to hash password: %w", err)
//	}
//	user.Password = hashedPassword
//
//	// Save the user to the database
//	err = u.repo.Create(user)
//	if err != nil {
//		return fmt.Errorf("failed to register user: %w", err)
//	}
//
//	// Send activation email or perform other registration tasks
//
//	return nil
//}
//
//func (u *userUseCaseImpl) Activation(token string) error {
//	// Implement activation logic using the token
//	// This might involve verifying the token against a database record
//	// and updating the user's activation status
//
//	return nil
//}
//
//func (u *userUseCaseImpl) ForgotPassword(email string) error {
//	// Implement logic to send a password reset email to the user
//	// This might involve generating a unique token, sending an email
//	// with a link containing the token, and updating the user's record
//	// with the token for verification later
//
//	return nil
//}
//
//func (u *userUseCaseImpl) Authenticate(email, password string) error {
//	// Retrieve
