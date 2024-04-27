package models

import (
	"time"
)

type Post struct {
	ID          int64
	CreatedAt   time.Time
	Name        string
	Description string
	AuthorID    int64
	Type        string
	Version     int64
}
