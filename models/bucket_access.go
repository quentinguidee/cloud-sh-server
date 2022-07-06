package models

type BucketAccess struct {
	Id         int
	BucketId   int
	UserId     int
	AccessType string
}

type AccessType int

const (
	Denied = iota
	ReadOnly
	Write
	Full
)
