package tokenRepository

import (
	"DiplomaV2/backend/internal/entity"
	"time"
)

type TokenRepository interface {
	New(userID int64, ttl time.Duration, scope string) (*entity.Token, error)
	insert(token *entity.Token) error
	DeleteAllForUser(scope string, userID int64) error
}
