package somesql

import (
	"strings"
)

// Fields constants
const (
	FieldData string = "data"
)

// Fields variables
var (
	ReservedFields = []string{"id", "created_at", "updated_at", "owner_id", "status", "type"}
)

// IsFieldMeta returns if field is a meta field
func IsFieldMeta(field string) bool {
	for _, f := range ReservedFields {
		if f == field {
			return true
		}
	}
	return false
}

// IsFieldData returns if field is a data field
func IsFieldData(field string) bool {
	return strings.HasPrefix(field, FieldData)
}

// GetFieldData returns data field with lang
func GetFieldData(lang string) string {
	return FieldData + "_" + lang
}
