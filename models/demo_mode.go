package models

type DemoMode struct {
	Enabled       bool   `json:"enabled" db:"enabled"`
	ResetInterval string `json:"reset_interval,omitempty" db:"reset_interval"`
}
