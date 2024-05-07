package entities

import "time"

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Username  string    `json:"username"`
	Telegram  string    `json:"telegram"`
	Discord   string    `json:"discord"`
	Email     string    `json:"email"`
	Skills    []string  `json:"skills"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	Plaintext *string
	Hash      []byte
}
