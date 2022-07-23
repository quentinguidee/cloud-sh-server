package models

type Server struct {
	ID              int `json:"id,omitempty" gorm:"primaryKey"`
	VersionMajor    int `json:"version_major,omitempty" gorm:"not null"`
	VersionMinor    int `json:"version_minor,omitempty" gorm:"not null"`
	VersionPatch    int `json:"version_patch,omitempty" gorm:"not null"`
	DatabaseVersion int `json:"database_version" gorm:"not null"`
}

type ServerVersion struct {
	Major int
	Minor int
	Patch int
}

func (server *Server) version() ServerVersion {
	return ServerVersion{
		Major: server.VersionMajor,
		Minor: server.VersionMinor,
		Patch: server.VersionPatch,
	}
}
