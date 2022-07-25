package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Node struct {
	UUID       string         `json:"uuid,omitempty" gorm:"primaryKey"`
	ParentUUID string         `json:"parent_uuid,omitempty" gorm:"default:NULL"`
	Parent     *Node          `json:"parent,omitempty" gorm:"foreignKey:ParentUUID"`
	BucketUUID uuid.UUID      `json:"bucket_uuid" gorm:"type:uuid;not null"`
	Name       string         `json:"name" gorm:"not null"`
	Type       string         `json:"type" gorm:"not null"`
	Mime       *string        `json:"mime,omitempty"`
	Size       *int64         `json:"size,omitempty"`
	NodeUsers  []NodeUser     `json:"users,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type NodeUser struct {
	UserID     int        `json:"user_id" gorm:"primaryKey"`
	NodeUUID   string     `json:"node_uuid" gorm:"primaryKey"`
	LastViewAt *time.Time `json:"last_view_at,omitempty"`
	EditedAt   *time.Time `json:"edited_at,omitempty"`
}
