package utils

import (
	"net/http"
	"self-hosted-cloud/server/database"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewTransaction(c *gin.Context) *sqlx.Tx {
	db := database.GetDatabaseFromContext(c)
	tx, err := db.Instance.BeginTxx(c, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}
	return tx
}

func ExecTransaction(c *gin.Context, tx *sqlx.Tx) {
	err := tx.Commit()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
