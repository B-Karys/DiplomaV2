package entities

import "time"

type (
	Post struct {
		ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
		CreatedAt   time.Time `json:"createdAt"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		AuthorID    int64     `json:"authorId"`
		Type        string    `json:"type"`
		Version     int64     `json:"version"`
	}
)
