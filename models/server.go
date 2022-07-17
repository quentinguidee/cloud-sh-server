package models

type Server struct {
	Id              int `json:"id,omitempty" db:"id"`
	VersionMajor    int `json:"version_major,omitempty" db:"version_major"`
	VersionMinor    int `json:"version_minor,omitempty" db:"version_minor"`
	VersionPatch    int `json:"version_patch,omitempty" db:"version_patch"`
	DatabaseVersion int `json:"database_version" db:"database_version"`
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
