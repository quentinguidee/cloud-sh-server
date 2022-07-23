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
	GithubAuths    []GithubAuth `json:"github_auths"`
	Nodes          []User       `json:"users" gorm:"many2many:node_users;"`
	Buckets        []Bucket     `json:"buckets" gorm:"many2many:node_users;"`
}

type GithubAuth struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
}
