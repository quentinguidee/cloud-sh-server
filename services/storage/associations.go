package storage

import (
	"errors"
	"net/http"
	. "self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

func CreateBucketNodeAssociation(tx *sqlx.Tx, fromNodeUuid string, toNodeUuid string) IServiceError {
	request := `
		INSERT INTO nodes_to_nodes(from_node, to_node)
		VALUES ($1, $2)
	`

	_, err := tx.Exec(request, fromNodeUuid, toNodeUuid)
	if err != nil {
		err = errors.New("error while creating bucket node association")
		return NewServiceError(http.StatusInternalServerError, err)
	}
	return nil
}
