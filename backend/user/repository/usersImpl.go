package repository

import (
	"DiplomaV2/backend/internal/database"
	"DiplomaV2/backend/user/models"
	"crypto/sha256"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type userRepository struct {
	DB database.Database
}

func NewUserRepository(db database.Database) UserRepository {
	return &userRepository{DB: db}
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

func (m *userRepository) Insert(user *models.User) error {
	result := m.DB.GetDb().Create(user).Scan(user)
	if result.Error != nil {
		switch {
		case result.Error.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return result.Error
		}
	}
	return nil
}

func (m *userRepository) GetByID(id int64) (*models.User, error) {
	// Initialize a User variable to store the result
	var user models.User

	// Query the database using GORM's First method
	result := m.DB.GetDb().First(&user, id)

	// Check for errors
	if result.Error != nil {
		// Handle specific errors, like record not found
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		// Handle other errors
		return nil, result.Error
	}

	return &user, nil
}

func (m *userRepository) GetByEmail(email string) (*models.User, error) {
	// Initialize a User variable to store the result
	var user models.User

	// Query the database using GORM's Where method
	result := m.DB.GetDb().Where("email = ?", email).First(&user)

	// Check for errors
	if result.Error != nil {
		// Handle specific errors, like record not found
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		// Handle other errors
		return nil, result.Error
	}
	return &user, nil
}

func (m *userRepository) Update(user *models.User) error {
	// Use GORM's Save method to update the user
	result := m.DB.GetDb().Save(user)
	// Check for errors
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *userRepository) GetForToken(tokenScope, tokenPlaintext string) (*models.User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Query the database using GORM's Where method
	var user models.User
	result := m.DB.GetDb().Joins("INNER JOIN tokens ON users.id = tokens.user_id").
		Where("tokens.hash = ?", tokenHash[:]).
		Where("tokens.scope = ?", tokenScope).
		Where("tokens.expiry > ?", time.Now()).
		First(&user)

	// Check for errors
	if result.Error != nil {
		// Handle specific errors, like record not found
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		// Handle other errors
		return nil, result.Error
	}

	return &user, nil
}

func (m *userRepository) Delete(id int64) error {
	// Construct the model instance with the ID set
	user := models.User{ID: id}

	// Delete the record using GORM's Delete method
	result := m.DB.GetDb().Delete(&user)

	// Check for errors
	if result.Error != nil {
		// Handle specific errors, like record not found
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		// Handle other errors
		return result.Error
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
