package entity

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt    time.Time      `gorm:"not null;default:current_timestamp" json:"created_at"`
	Name         string         `gorm:"not null" json:"name"`
	Surname      string         `json:"surname"`
	Username     string         `gorm:"unique; not null" json:"username"`
	Telegram     string         `json:"telegram"`
	Discord      string         `json:"discord"`
	Email        string         `gorm:"type:citext;unique;not null" json:"email"`
	Skills       pq.StringArray `gorm:"type:text[]" json:"skills"`
	Password     password       `gorm:"embedded;embeddedPrefix:password_" json:"-"`
	ProfileImage string         `json:"profileImage"`
	Activated    bool           `gorm:"default:false;not null" json:"activated"`
	Version      int            `gorm:"not null;default:1" json:"-"`
	Posts        []Post         `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;" json:"posts"`
	Tokens       []Token        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"tokens"`
}

type Token struct {
	ID        uint      `gorm:"primaryKey"`
	Plaintext string    `gorm:"-"`
	Hash      []byte    `gorm:"not null" json:"-"`
	UserID    int64     `gorm:"foreignKey:user_id;not null" json:"-"`
	Expiry    time.Time `gorm:"not null" json:"expiry"`
	Scope     string    `gorm:"not null" json:"-"`
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
