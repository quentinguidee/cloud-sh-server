package models

type Bucket struct {
	Id   int    `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
	Type string `json:"type,omitempty" db:"type"`
}
