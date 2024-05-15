package repository

import (
	"DiplomaV2/token/models"
	"time"
)

type TokenRepository interface {
	New(userID int64, ttl time.Duration, scope string) (*models.Token, error)
	Insert(token *models.Token) error
	DeleteAllForUser(scope string, userID int64) error
}
