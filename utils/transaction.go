package utils

import (
	"database/sql"
	"net/http"
	"self-hosted-cloud/server/database"

	"github.com/gin-gonic/gin"
)

func NewTransaction(c *gin.Context) *sql.Tx {
	db := database.GetDatabaseFromContext(c)
	tx, err := db.Instance.BeginTx(c, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return nil
	}
	return tx
}

func ExecTransaction(c *gin.Context, tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
