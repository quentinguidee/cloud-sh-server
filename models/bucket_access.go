package models

type BucketAccess struct {
	Id         int    `json:"id,omitempty" db:"id"`
	BucketId   int    `json:"bucket_id,omitempty" db:"bucket_id"`
	UserId     int    `json:"user_id,omitempty" db:"user_id"`
	AccessType string `json:"access_type,omitempty" db:"access_type"`
}

type AccessType int

const (
	Denied = iota
	ReadOnly
	Write
	Full
)
