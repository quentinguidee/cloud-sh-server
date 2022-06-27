package storage

import (
	"fmt"
	"os"
)

type Bucket struct {
	Id       int
	Name     string
	RootNode int
	Type     string
}

// Create creates the bucket in the localstorage.
func (bucket Bucket) Create() error {
	err := os.MkdirAll(fmt.Sprintf("%s/buckets/%d", os.Getenv("DATA_PATH"), bucket.Id), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
