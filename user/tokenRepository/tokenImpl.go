package tokenRepository

import (
	"DiplomaV2/database"
	"DiplomaV2/user/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"github.com/pkg/errors"
	"time"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password-reset"
)

type tokenRepository struct {
	DB database.Database
}

func (t *tokenRepository) New(userID int64, ttl time.Duration, scope string) (*models.Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.insert(token)
	return token, err
}

func (t *tokenRepository) insert(token *models.Token) error {
	if err := t.DB.GetDb().Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (t *tokenRepository) DeleteAllForUser(scope string, userID int64) error {
	if err := t.DB.GetDb().Where("scope = ? AND user_id = ?", scope, userID).Delete(&models.Token{}).Error; err != nil {
		return err
	}
	return nil
}
func (t *tokenRepository) DeleteToken(userID int64, tokenString string) error {
	var token models.Token

	// Find the token for the user
	if err := t.DB.GetDb().Where("user_id = ?", userID).First(&token).Error; err != nil {
		return err
	}

	// Compare the token string with the stored plaintext token
	if tokenString != token.Plaintext {
		return errors.New("token does not match the provided user")
	}

	// Delete the token
	if err := t.DB.GetDb().Delete(&token).Error; err != nil {
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
