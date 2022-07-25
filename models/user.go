package models

import (
	"time"
)

type User struct {
	ID             int          `json:"id" gorm:"primaryKey"`
	Username       string       `json:"username" gorm:"unique,not null"`
	Name           string       `json:"name" gorm:"not null"`
	ProfilePicture *string      `json:"profile_picture,omitempty"`
	Role           *string      `json:"role,omitempty"`
	CreatedAt      time.Time    `json:"created_at,omitempty" gorm:"not null"`
	Sessions       []Session    `json:"sessions,omitempty"`
	GithubUsers    []GithubUser `json:"github_users,omitempty"`
	Nodes          []User       `json:"users,omitempty" gorm:"many2many:node_users;"`
	Buckets        []Bucket     `json:"buckets,omitempty" gorm:"many2many:user_buckets;"`
}

type GithubUser struct {
	Username string `json:"username" gorm:"primaryKey"`
	UserID   int    `json:"user_id" gorm:"primaryKey"`
	User     User   `json:"user" gorm:"foreignKey:UserID"`
}
