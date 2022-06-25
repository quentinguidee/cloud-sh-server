package models

type User struct {
	Id             int    `json:"id,omitempty"`
	Username       string `json:"username,omitempty"`
	Name           string `json:"name,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

type GithubUser struct {
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	Login     string `json:"login,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}
