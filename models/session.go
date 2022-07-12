package models

type Session struct {
	Id     int    `json:"id,omitempty" db:"id"`
	UserId int    `json:"user_id,omitempty" db:"user_id"`
	Token  string `json:"token,omitempty" db:"token"`
}
