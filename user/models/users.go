package models

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Username  string    `json:"username"`
	Telegram  string    `json:"telegram"`
	Discord   string    `json:"discord"`
	Email     string    `json:"email"`
	Skills    []string  `gorm:"type:text[]" json:"skills"`
	Password  password  `gorm:"embedded;embeddedPrefix:password_" json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type Token struct {
	Plaintext string    `gorm:"-"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type password struct {
	Plaintext *string `gorm:"-"`
	Hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
