package models

import (
	"github.com/lib/pq"
	"time"
)

type Post struct {
	ID          int64          `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt   time.Time      `gorm:"not null;default:current_timestamp" json:"createdAt" `
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"" json:"description"`
	AuthorID    int64          `gorm:"not null" json:"authorId"`
	Type        string         `gorm:"not null" json:"type"`
	Skills      pq.StringArray `gorm:"type:text[];" json:"skills"`
	Version     int            `gorm:"not null; default:1" json:"-"`
}
