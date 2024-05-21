package tokenRepository

import (
	"DiplomaV2/user/models"
	"time"
)

type TokenRepository interface {
	New(userID int64, ttl time.Duration, scope string) (*models.Token, error)
	insert(token *models.Token) error
	DeleteToken(userId int64, token string) error
	DeleteAllForUser(scope string, userID int64) error
}
