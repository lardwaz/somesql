package somesql

import (
	"sort"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Fields constants
const (
	FieldID        string = "id"
	FieldCreatedAt string = "created_at"
	FieldUpdatedAt string = "updated_at"
	FieldOwnerID   string = "owner_id"
	FieldStatus    string = "status"
	FieldType      string = "type"
	FieldData      string = "data"
	FieldRelations string = "relations"
)

// Fields variables
var (
	MetaFieldsList = []string{FieldID, FieldCreatedAt, FieldUpdatedAt, FieldOwnerID, FieldStatus, FieldType}
	FieldsList     = append(MetaFieldsList, FieldData, FieldRelations)
)

// Fields represents Top Level Fields
type Fields map[string]interface{}

// NewFields return new Fields
func NewFields() Fields {
	fields := make(Fields)

	return fields
}

// UseDefaults sets the default values for Fields
func (f Fields) UseDefaults() Fields {
	// Default fields values
	f[FieldID] = uuid.NewV4().String()
	f[FieldCreatedAt] = time.Now()
	f[FieldUpdatedAt] = time.Now()
	f[FieldOwnerID] = uuid.Nil.String()
	f[FieldStatus] = ""
	f[FieldType] = ""
	f[FieldData] = make(map[string]interface{})
	f[FieldRelations] = make(map[string]interface{})

	return f
}

// ID is a setter for ID in Fields
func (f Fields) ID(id string) Fields {
	f.Set(FieldID, id)

	return f
}

// CreatedAt is a setter for CreatedAt in Fields
func (f Fields) CreatedAt(t time.Time) Fields {
	f.Set(FieldCreatedAt, t)

	return f
}

// UpdatedAt is a setter for UpdatedAt in Fields
func (f Fields) UpdatedAt(t time.Time) Fields {
	f.Set(FieldUpdatedAt, t)

	return f
}

// OwnerID is a setter for OwnerID in Fields
func (f Fields) OwnerID(id string) Fields {
	f.Set(FieldOwnerID, id)

	return f
}

// Status is a setter for Status in Fields
func (f Fields) Status(s string) Fields {
	f.Set(FieldStatus, s)

	return f
}

// Type is a setter for Type in Fields
func (f Fields) Type(s string) Fields {
	f.Set(FieldType, s)

	return f
}

// Set assigns a new value to fields
// Dot-seperated field name is treated as inner field of JSONB field (1 level only)
// i.e data.author = data->>author
func (f Fields) Set(field string, value interface{}) Fields {
	if innerField := GetInnerDataField(field); innerField != "" {
		innerValue, ok := f[FieldData].(map[string]interface{})
		if !ok {
			f[FieldData] = make(map[string]interface{})
			innerValue = f[FieldData].(map[string]interface{})
		}
		innerValue[innerField] = value
	} else if innerField := GetInnerRelationsField(field); innerField != "" {
		innerValue, ok := f[FieldRelations].(map[string]interface{})
		if !ok {
			f[FieldRelations] = make(map[string]interface{})
			innerValue = f[FieldRelations].(map[string]interface{})
		}
		innerValue[innerField] = value
	}

	for _, ff := range FieldsList {
		if ff == field {
			f[field] = value
		}
	}

	return f
}

// List returns fields and values
func (f Fields) List() ([]string, []interface{}) {
	fields := make([]string, 0)
	values := make([]interface{}, 0)

	// Sort the map by keys first
	var keys []string
	for k := range f {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// First add top level fields in order of FieldsList
	for _, field := range FieldsList {
		if v, ok := f[field]; ok {
			fields = append(fields, field)
			values = append(values, v)
		}
	}

	// Then add inner fields
	for field := range f {
		if strings.Contains(field, ".") {
			fields = append(fields, field)
			values = append(values, f[field])
		}
	}

	return fields, values
}

// IsFieldMeta returns true if field is a meta field
func IsFieldMeta(field string) bool {
	for _, f := range FieldsList {
		if f == field {
			return true
		}
	}
	return false
}

// IsFieldData returns true if field is a data field
func IsFieldData(field string) bool {
	return field == FieldData
}

// GetInnerDataField returns the inner data field
func GetInnerDataField(field string) string {
	if strings.Count(field, ".") == 1 {
		parts := strings.Split(field, ".")
		if IsFieldData(parts[0]) {
			return parts[1]
		}
	}
	return ""
}

// IsFieldRelations returns true if field is a data field
func IsFieldRelations(field string) bool {
	return field == FieldRelations
}

// GetInnerRelationsField returns the inner relations field
func GetInnerRelationsField(field string) string {
	if strings.Count(field, ".") == 1 {
		parts := strings.Split(field, ".")
		if IsFieldRelations(parts[0]) {
			return parts[1]
		}
	}
	return ""
}

// GetFieldData returns data field with lang
func GetFieldData(lang string) string {
	return FieldData + "_" + lang
}
