package models

type Bucket struct {
	Id       int    `json:"id,omitempty" db:"id"`
	Name     string `json:"name,omitempty" db:"name"`
	RootNode string `json:"root_node,omitempty" db:"root_node"`
	Type     string `json:"type,omitempty" db:"type"`
}
