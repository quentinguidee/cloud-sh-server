package models

type User struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
}

type GithubUser struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Login string `json:"login,omitempty"`
}
