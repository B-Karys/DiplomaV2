package repository

import (
	"DiplomaV2/database"
	"DiplomaV2/token/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"DiplomaV2/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password-reset"
)

type tokenRepository struct {
	DB database.Database
}

func (t tokenRepository) New(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	return token, err
}

func (t tokenRepository) Insert(token *models.Token) error {
	if err := t.DB.GetDb().Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (t tokenRepository) DeleteAllForUser(scope string, userID int64) error {
	if err := t.DB.GetDb().Where("scope = ? AND user_id = ?", scope, userID).Delete(&models.Token{}).Error; err != nil {
		return err
	}
	return nil
}

func NewTokenRepository(db database.Database) TokenRepository {
	return &tokenRepository{DB: db}
}

func generateToken(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	token := &models.Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
