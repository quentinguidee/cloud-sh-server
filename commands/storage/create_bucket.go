package commands

import (
	"fmt"
	"net/http"
	"os"
	. "self-hosted-cloud/server/commands"
	. "self-hosted-cloud/server/database"
	. "self-hosted-cloud/server/models/storage"
)

type CreateBucketTransaction struct {
	Bucket   *Bucket
	Database Database
	UserId   int
}

func (c CreateBucketTransaction) Try() ICommandError {
	node := Node{
		Filename: "/",
		Filetype: "directory",
	}

	commands := []Command{
		CreateBucketInDatabaseCommand{
			Bucket:   c.Bucket,
			Database: c.Database,
			UserId:   c.UserId,
		},
		CreateBucketNodeCommand{
			Node:     &node,
			Bucket:   c.Bucket,
			Database: c.Database,
		},
		UpdateBucketRootNodeCommand{
			Bucket:   c.Bucket,
			Database: c.Database,
			Node:     &node,
		},
		CreateBucketAccess{
			Bucket:   c.Bucket,
			Database: c.Database,
			UserId:   c.UserId,
		},
		CreateBucketInFileSystemCommand{
			Bucket: c.Bucket,
		},
	}
	return NewTransaction(commands).Try()
}

type CreateBucketInFileSystemCommand struct {
	Bucket *Bucket
}

func (c CreateBucketInFileSystemCommand) Run() ICommandError {
	err := os.MkdirAll(fmt.Sprintf("%s/buckets/%d", os.Getenv("DATA_PATH"), c.Bucket.Id), os.ModePerm)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateBucketInFileSystemCommand) Revert() ICommandError {
	err := os.RemoveAll(fmt.Sprintf("%s/buckets/%d", os.Getenv("DATA_PATH"), c.Bucket.Id))
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

type CreateBucketInDatabaseCommand struct {
	Bucket   *Bucket
	Database Database
	UserId   int
}

func (c CreateBucketInDatabaseCommand) Run() ICommandError {
	request := "INSERT INTO buckets(name, type) VALUES (?, ?) RETURNING id"

	err := c.Database.Instance.QueryRow(request, c.Bucket.Name, c.Bucket.Type).Scan(&c.Bucket.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}

func (c CreateBucketInDatabaseCommand) Revert() ICommandError {
	request := "DELETE FROM buckets WHERE id = ?"

	_, err := c.Database.Instance.Exec(request, c.Bucket.Id)
	if err != nil {
		return NewError(http.StatusInternalServerError, err)
	}
	return nil
}
