package models

type User struct {
	Id             int    `json:"id,omitempty" db:"id"`
	Username       string `json:"username,omitempty" db:"username"`
	Name           string `json:"name,omitempty" db:"name"`
	ProfilePicture string `json:"profile_picture,omitempty" db:"profile_picture"`
}

type GithubUser struct {
	Email     string `json:"email,omitempty" db:"email"`
	Name      string `json:"name,omitempty" db:"name"`
	Login     string `json:"login,omitempty" db:"login"`
	AvatarUrl string `json:"avatar_url,omitempty" db:"avatar_url"`
}
