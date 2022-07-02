package storage

type Node struct {
	Uuid     string `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	Filetype string `json:"filetype,omitempty"`
	BucketId int    `json:"bucket_id,omitempty"`
}
