package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullableTime struct {
	sql.NullTime
}

func NewNullableTime(time time.Time) NullableTime {
	return NullableTime{sql.NullTime{
		Time:  time,
		Valid: true,
	}}
}

func (t NullableTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}
