package storage

import (
	"os"
	"path/filepath"
	"strconv"
)

type Node struct {
	Id       int    `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	Filetype string `json:"filetype,omitempty"`
	BucketId int    `json:"bucket_id,omitempty"`
}

func (node Node) Delete(path string) error {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	path = filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(node.BucketId), path)
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
