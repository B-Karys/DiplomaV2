package repository

import (
	"DiplomaV2/backend/internal/database"
	"DiplomaV2/backend/internal/entity"
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

func (m *userRepository) Insert(user *entity.User) error {
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

func (m *userRepository) GetByID(id int64) (*entity.User, error) {
	var user entity.User

	result := m.DB.GetDb().First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}

	return &user, nil
}

func (m *userRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User

	result := m.DB.GetDb().Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &user, nil
}

func (m *userRepository) Update(user *entity.User) error {
	result := m.DB.GetDb().Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *userRepository) GetForToken(tokenScope, tokenPlaintext string) (*entity.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	var user entity.User
	result := m.DB.GetDb().Joins("INNER JOIN tokens ON users.id = tokens.user_id").
		Where("tokens.hash = ?", tokenHash[:]).
		Where("tokens.scope = ?", tokenScope).
		Where("tokens.expiry > ?", time.Now()).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}

	return &user, nil
}

func (m *userRepository) Delete(id int64) error {
	user := entity.User{ID: id}

	result := m.DB.GetDb().Delete(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
