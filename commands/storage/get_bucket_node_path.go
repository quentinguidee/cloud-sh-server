package commands

import (
	"os"
	"path/filepath"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
	"strconv"
)

type GetBucketNodePathCommand struct {
	Database     Database
	Path         string
	Bucket       *Bucket
	CompletePath *string
}

func (c GetBucketNodePathCommand) Run() ICommandError {
	// TODO before release: Check if user has permission to get this link

	*c.CompletePath = filepath.Join(os.Getenv("DATA_PATH"), "buckets", strconv.Itoa(c.Bucket.Id), c.Path)

	return nil
}

func (c GetBucketNodePathCommand) Revert() ICommandError {
	return nil
}
