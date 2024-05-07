package models

type User struct {
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	Username string   `json:"username"`
	Telegram string   `json:"telegram"`
	Discord  string   `json:"discord"`
	Email    string   `json:"email"`
	Skills   []string `json:"skills"`
}
