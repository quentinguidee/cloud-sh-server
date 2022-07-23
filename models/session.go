package models

type Session struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	UserID int    `json:"user_id" gorm:"not null"`
	Token  string `json:"token" gorm:"unique,not null"`
}
