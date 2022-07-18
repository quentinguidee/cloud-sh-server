package types

import (
	"database/sql"
	"encoding/json"
)

type NullableInt64 struct {
	sql.NullInt64
}

func NewNullableInt64(value int64) NullableInt64 {
	return NullableInt64{sql.NullInt64{
		Int64: value,
		Valid: true,
	}}
}

func (i NullableInt64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Int64)
}
