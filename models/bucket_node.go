package models

type Node struct {
	Uuid     string `json:"uuid,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	BucketId int    `json:"bucket_id,omitempty"`
}
