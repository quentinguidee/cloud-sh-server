package models

type DemoMode struct {
	Enabled       bool   `json:"enabled"`
	ResetInterval string `json:"reset_interval,omitempty"`
}
