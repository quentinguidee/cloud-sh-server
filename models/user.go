package models

import (
	. "self-hosted-cloud/server/models/types"
)

type User struct {
	Id             int            `json:"id,omitempty" db:"id"`
	Username       string         `json:"username,omitempty" db:"username"`
	Name           string         `json:"name,omitempty" db:"name"`
	ProfilePicture NullableString `json:"profile_picture,omitempty" db:"profile_picture"`
	Role           NullableString `json:"role,omitempty" db:"role"`
	CreationDate   NullableTime   `json:"creation_date,omitempty" db:"creation_date"`
}

type GithubUser struct {
	Email     string `json:"email,omitempty" db:"email"`
	Name      string `json:"name,omitempty" db:"name"`
	Login     string `json:"login,omitempty" db:"login"`
	AvatarUrl string `json:"avatar_url,omitempty" db:"avatar_url"`
}
