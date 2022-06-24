package models

type User struct {
	Id       int
	Username string
	Name     string
}

type GithubUser struct {
	Email string
	Name  string
	Login string
}
