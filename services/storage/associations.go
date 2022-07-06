package storage

import (
	"database/sql"
	"errors"
	"net/http"
	. "self-hosted-cloud/server/services"
)

func CreateBucketNodeAssociation(tx *sql.Tx, fromNodeUuid string, toNodeUuid string) IServiceError {
	request := `
		INSERT INTO buckets_nodes_associations(from_node, to_node)
		VALUES (?, ?)
	`

	_, err := tx.Exec(request, fromNodeUuid, toNodeUuid)
	if err != nil {
		err = errors.New("error while creating bucket node association")
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}
