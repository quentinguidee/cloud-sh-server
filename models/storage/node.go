package storage

type Node struct {
	Id               int    `json:"id,omitempty"`
	Filename         string `json:"filename,omitempty"`
	Filetype         string `json:"filetype,omitempty"`
	InternalFilename string `json:"internal_filename,omitempty"`
}
