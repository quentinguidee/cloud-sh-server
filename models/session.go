package models

type Session struct {
	Id     int    `json:"id,omitempty"`
	UserId int    `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}
