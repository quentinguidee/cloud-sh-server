package models

type Node struct {
	Uuid     string `json:"uuid,omitempty" db:"uuid"`
	Name     string `json:"name,omitempty" db:"name"`
	Type     string `json:"type,omitempty" db:"type"`
	Mime     string `json:"mime,omitempty" db:"mime"`
	Size     int64  `json:"size,omitempty" db:"size"`
	BucketId int    `json:"bucket_id,omitempty" db:"bucket_id"`
}
