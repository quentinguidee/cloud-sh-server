package commands

import (
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
	"strconv"
)

type GetBucketNodePath struct {
	Database     Database
	Path         string
	Bucket       *Bucket
	CompletePath *string
}

func (c GetBucketNodePath) Run() ICommandError {
	// TODO before release: Check if user has permission to get this link

	*c.CompletePath = filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(c.Bucket.Id), c.Path)

	return nil
}

func (c GetBucketNodePath) Revert() ICommandError {
	return nil
}
