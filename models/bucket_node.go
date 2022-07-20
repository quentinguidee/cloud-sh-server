package models

import . "self-hosted-cloud/server/models/types"

type Node struct {
	Uuid       string         `json:"uuid,omitempty" db:"uuid"`
	ParentUuid NullableString `json:"parent_uuid" db:"parent_uuid"`
	Name       string         `json:"name,omitempty" db:"name"`
	Type       string         `json:"type,omitempty" db:"type"`
	Mime       NullableString `json:"mime,omitempty" db:"mime"`
	Size       NullableInt64  `json:"size,omitempty" db:"size"`
}
