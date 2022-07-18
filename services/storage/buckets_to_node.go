package storage

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/models"
	. "self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

func CreateBucketToNodeAssociation(tx *sqlx.Tx, bucketId int, nodeUuid string) IServiceError {
	request := `
		INSERT INTO buckets_to_node(bucket_id, node_id)
		VALUES ($1, $2)
	`

	_, err := tx.Exec(request, bucketId, nodeUuid)
	if err != nil {
		err = errors.New("error while creating bucket node association")
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}

func GetBucketFromNode(tx *sqlx.Tx, nodeUuid string) (Bucket, IServiceError) {
	request := `
		SELECT buckets.*
		FROM buckets_to_node, buckets
		WHERE node_id = $1
		  AND bucket_id = buckets.id
	`

	var bucket Bucket
	err := tx.Get(&bucket, request, nodeUuid)
	if err != nil {
		return Bucket{}, NewServiceError(http.StatusInternalServerError, err)
	}
	return bucket, nil
}
