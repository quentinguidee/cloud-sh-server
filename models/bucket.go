package models

import "github.com/google/uuid"

type Bucket struct {
	UUID     uuid.UUID `json:"uuid" gorm:"type:uuid;primaryKey"`
	Name     string    `json:"name" gorm:"not null"`
	Type     string    `json:"type" gorm:"not null"`
	Size     int64     `json:"size" gorm:"not null;default:0"`
	RootNode *Node     `json:"root_node" gorm:"-:all"`
	MaxSize  *int64    `json:"max_size"`
	Users    []User    `json:"users" gorm:"many2many:user_buckets;"`
}

type UserBucket struct {
	BucketUUID uuid.UUID `json:"bucket_uuid" gorm:"not null;type:uuid"`
	Bucket     Bucket    `json:"bucket" gorm:"foreignKey:BucketUUID;not null"`
	UserID     int       `json:"user_id" gorm:"not null"`
	User       User      `json:"user" gorm:"foreignKey:UserID;not null"`
	AccessType string    `json:"access_type" gorm:"not null"`
}

type AccessType int

const (
	Denied = iota
	ReadOnly
	Write
	Full
)
