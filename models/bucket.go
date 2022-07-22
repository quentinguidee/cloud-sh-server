package models

import "self-hosted-cloud/server/models/types"

type Bucket struct {
	Id      int                 `json:"id,omitempty" db:"id"`
	Name    string              `json:"name,omitempty" db:"name"`
	Type    string              `json:"type,omitempty" db:"type"`
	Size    int64               `json:"size" db:"size"`
	MaxSize types.NullableInt64 `json:"max_size" db:"max_size"`
}
