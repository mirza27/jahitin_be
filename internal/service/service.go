package service

import (
	"database/sql"
	"time"
)

// helper function
func TimeOrZero(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}

	return *value
}

func NullStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func NullInt64Value(i sql.NullInt64) int64 {
	if i.Valid {
		return i.Int64
	}
	return 0
}

func nullInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *value, Valid: true}
}

func nullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *value, Valid: true}
}

func NullTimeValue(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}
