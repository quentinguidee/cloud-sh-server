package types

import (
	"database/sql"
	"encoding/json"
)

type NullableString struct {
	sql.NullString
}

func NewNullString() NullableString {
	return NullableString{sql.NullString{
		String: "",
		Valid:  false,
	}}
}

func NewNullableString(string string) NullableString {
	return NullableString{sql.NullString{
		String: string,
		Valid:  true,
	}}
}

func (s NullableString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}
