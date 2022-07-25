package models

import (
	"time"
)

type User struct {
	ID             int          `json:"id" gorm:"primaryKey"`
	Username       string       `json:"username" gorm:"unique,not null"`
	Name           string       `json:"name" gorm:"not null"`
	ProfilePicture *string      `json:"profile_picture"`
	Role           *string      `json:"role"`
	CreatedAt      time.Time    `json:"created_at" gorm:"not null"`
	Sessions       []Session    `json:"sessions"`
	GithubUsers    []GithubUser `json:"github_users"`
	Nodes          []User       `json:"users" gorm:"many2many:node_users;"`
	Buckets        []Bucket     `json:"buckets" gorm:"many2many:user_buckets;"`
}

type GithubUser struct {
	Username string `json:"username" gorm:"primaryKey"`
	UserID   int    `json:"user_id" gorm:"primaryKey"`
	User     User   `json:"user" gorm:"foreignKey:UserID"`
}
