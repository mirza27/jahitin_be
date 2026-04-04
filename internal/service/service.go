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

func NullTimeValue(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}
