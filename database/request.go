package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"self-hosted-cloud/server/services"

	"github.com/jmoiron/sqlx"
)

type Request struct {
	tx    *sqlx.Tx
	query string
}

func NewRequest(tx *sqlx.Tx, query string) *Request {
	return &Request{tx, query}
}

func (r *Request) QueryRow(args ...any) *Row {
	row := r.tx.QueryRowx(r.query, args...)
	return &Row{row}
}

func (r *Request) Exec(args ...any) *ExecResult {
	result, err := r.tx.Exec(r.query, args...)
	return &ExecResult{result, err}
}

func (r *Request) Get(dest interface{}, args ...any) *RequestResult {
	err := r.tx.Get(dest, r.query, args...)
	return &RequestResult{err}
}

func (r *Request) Select(dest interface{}, args ...any) *RequestResult {
	err := r.tx.Select(dest, r.query, args...)
	return &RequestResult{err}
}

type Row struct {
	row *sqlx.Row
}

func (r *Row) Scan(args ...any) *RequestResult {
	err := r.row.Scan(args...)
	return &RequestResult{err}
}

type ExecResult struct {
	result sql.Result
	err    error
}

type RequestResult struct {
	err error
}

func (r *RequestResult) OnError(err string) services.IServiceError {
	return onError(r.err, err)
}

func (r *ExecResult) OnError(err string) (sql.Result, services.IServiceError) {
	return r.result, onError(r.err, err)
}

func onError(requestError error, err string) services.IServiceError {
	if requestError == nil {
		return nil
	}

	serviceError := fmt.Errorf("%s: %s", err, requestError)
	if requestError == sql.ErrNoRows {
		return services.NewServiceError(http.StatusNotFound, serviceError)
	}
	return services.NewServiceError(http.StatusInternalServerError, serviceError)
}
