package models

type Bucket struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	Type     string `json:"type" gorm:"not null"`
	Size     int64  `json:"size" gorm:"not null,default:0"`
	RootNode *Node  `json:"root_node" gorm:"-:all"`
	MaxSize  *int64 `json:"max_size"`
	Users    []User `json:"users" gorm:"many2many:bucket_users;"`
}

type BucketUser struct {
	BucketID   int    `json:"bucket_id" gorm:"not null"`
	UserID     int    `json:"user_id" gorm:"not null"`
	AccessType string `json:"access_type" gorm:"not null"`
}

type AccessType int

const (
	Denied = iota
	ReadOnly
	Write
	Full
)
