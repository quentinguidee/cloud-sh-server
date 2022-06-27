package storage

import (
	"errors"
	"fmt"
	"os"
)

type Node struct {
	Id       int    `json:"id,omitempty"`
	Filename string `json:"filename,omitempty"`
	Filetype string `json:"filetype,omitempty"`
	BucketId int    `json:"bucket_id,omitempty"`
}

// Create creates the node in the localstorage.
func (node Node) Create(path string) error {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if len(path) > 0 {
		path += "/"
	}

	var err error

	switch node.Filetype {
	case "directory":
		err = os.Mkdir(fmt.Sprintf("%s/buckets/%d/%s%s", os.Getenv("DATA_PATH"), node.BucketId, path, node.Filename), os.ModePerm)
	default:
		err = errors.New("this filetype is not supported")
	}

	if err != nil {
		return err
	}

	return nil
}
