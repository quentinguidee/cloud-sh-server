package models

import (
	"gorm.io/gorm"
	"time"
)

type Node struct {
	UUID       string         `json:"uuid,omitempty" gorm:"primaryKey"`
	ParentUUID string         `json:"parent_uuid" gorm:"default:NULL"`
	Parent     *Node          `json:"parent" gorm:"foreignKey:ParentUUID"`
	BucketID   int            `json:"bucket_id" gorm:"not null"`
	Name       string         `json:"name" gorm:"not null"`
	Type       string         `json:"type" gorm:"not null"`
	Mime       *string        `json:"mime"`
	Size       *int64         `json:"size"`
	NodeUsers  []NodeUser     `json:"users"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type NodeUser struct {
	UserID     int        `json:"user_id" gorm:"primaryKey"`
	NodeUUID   string     `json:"node_uuid" gorm:"primaryKey"`
	LastViewAt *time.Time `json:"last_view_at"`
	EditedAt   *time.Time `json:"edited_at"`
}
