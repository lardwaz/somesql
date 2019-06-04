package somesql

import (
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

	JSONBArrSet uint8 = iota
	JSONBArrAdd
	JSONBArrRemove
)

// Fields variables
var (
	MetaFieldsList = []string{FieldID, FieldCreatedAt, FieldUpdatedAt, FieldOwnerID, FieldStatus, FieldType}
	FieldsList     = append(MetaFieldsList, FieldData, FieldRelations)
)

// JSONBField represents information about a single JSONB field
type JSONBField struct {
	Value  interface{}
	Action uint8
}

// JSONBFields represents multiple jsonb fields
type JSONBFields struct {
	data map[string]JSONBField
	keys []string // track order of insertion TODO: check other patterns?
}

// NewJSONBFields returns a new JSONBFields
func NewJSONBFields() JSONBFields {
	return JSONBFields{
		data: make(map[string]JSONBField),
		keys: make([]string, 0),
	}
}

// Add adds a new key-value to data
func (j *JSONBFields) Add(field string, value interface{}, action ...uint8) {
	act := JSONBArrSet
	if len(action) == 1 {
		act = action[0]
	}
	j.data[field] = JSONBField{Value: value, Action: act}
	j.keys = append(j.keys, field)
}

// GetOrderedList returns ordered list of fields name and value in insertion order
func (j JSONBFields) GetOrderedList() ([]string, []interface{}, []uint8) {
	var (
		values  []interface{}
		actions []uint8
	)
	for _, k := range j.keys {
		jsonbField := j.data[k]
		values = append(values, jsonbField.Value)
		actions = append(actions, jsonbField.Action)
	}

	return j.keys, values, actions
}

// Values returns the inner map
func (j JSONBFields) Values() map[string]interface{} {
	values := make(map[string]interface{})
	for f, v := range j.data {
		values[f] = v.Value
	}
	return values
}

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
	f[FieldData] = NewJSONBFields()
	f[FieldRelations] = NewJSONBFields()

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
	if IsFieldMeta(field) || IsWholeFieldData(field) || IsWholeFieldRelations(field) {
		f[field] = value
	} else if innerField, ok := GetInnerField(FieldData, field); ok {
		jsonbFields, ok := f[FieldData].(JSONBFields)
		if !ok { // if not jsonbfields, make it
			jsonbFields = NewJSONBFields()
		}
		jsonbFields.Add(innerField, value)
		f[FieldData] = jsonbFields
	} else if innerField, ok := GetInnerField(FieldRelations, field); ok {
		jsonbFields, ok := f[FieldRelations].(JSONBFields)
		if !ok { // if not jsonbfields, make it
			jsonbFields = NewJSONBFields()
		}
		jsonbFields.Add(innerField, value)
		f[FieldRelations] = jsonbFields
	}

	return f
}

// AddToArray adds data to an array jsonb field
func (f Fields) AddToArray(field string, value interface{}) Fields {
	if innerField, ok := GetInnerField(FieldData, field); ok {
		jsonbFields, ok := f[FieldData].(JSONBFields)
		if !ok { // if not jsonbfields, make it
			jsonbFields = NewJSONBFields()
		}
		jsonbFields.Add(innerField, value, JSONBArrAdd)
		f[FieldData] = jsonbFields
	} else if innerField, ok := GetInnerField(FieldRelations, field); ok {
		jsonbFields, ok := f[FieldRelations].(JSONBFields)
		if !ok {
			jsonbFields = NewJSONBFields()
		}
		jsonbFields.Add(innerField, value, JSONBArrAdd)
		f[FieldRelations] = jsonbFields
	}

	return f
}

// RemoveFromArray remove data from an array jsonb field
func (f Fields) RemoveFromArray(field string, value interface{}) Fields {
	if innerField, ok := GetInnerField(FieldData, field); ok {
		jsonbFields, ok := f[FieldData].(JSONBFields)
		if !ok { // if not jsonbfields, make it
			jsonbFields = NewJSONBFields()
		}
		jsonbFields.Add(innerField, value, JSONBArrRemove)
		f[FieldData] = jsonbFields
	} else if innerField, ok := GetInnerField(FieldRelations, field); ok {
		jsonbFields, ok := f[FieldRelations].(JSONBFields)
		if !ok {
			jsonbFields = NewJSONBFields()
		}
		jsonbFields.Add(innerField, value, JSONBArrRemove)
		f[FieldRelations] = jsonbFields
	}

	return f
}

// List returns fields and values
func (f Fields) List() ([]string, []interface{}) {
	fields := make([]string, 0)
	values := make([]interface{}, 0)

	// Meta Fields
	for _, field := range MetaFieldsList {
		if v, ok := f[field]; ok {
			fields = append(fields, field)
			values = append(values, v)
		}
	}

	// Data Fields
	if dataField, ok := f[FieldData].(JSONBFields); ok {
		fields = append(fields, FieldData)
		values = append(values, dataField)
	}

	// Relations Fields
	if relationsField, ok := f[FieldRelations].(JSONBFields); ok {
		fields = append(fields, FieldRelations)
		values = append(values, relationsField)
	}

	return fields, values
}

// IsFieldMeta returns true if field is a meta field
func IsFieldMeta(field string) bool {
	for _, f := range MetaFieldsList {
		if f == field {
			return true
		}
	}
	return false
}

// IsWholeFieldData returns true if field is a data field
func IsWholeFieldData(field string) bool {
	return field == FieldData
}

// IsWholeFieldRelations returns true if field is a data field
func IsWholeFieldRelations(field string) bool {
	return field == FieldRelations
}

// GetLangFieldData returns data field with lang
func GetLangFieldData(lang string) string {
	return FieldData + "_" + lang
}

// GetInnerField returns the inner data field
func GetInnerField(parent, field string) (string, bool) {
	if strings.Count(field, ".") == 1 {
		parts := strings.Split(field, ".")
		if parts[0] != parent && parts[1] != "" {
			return "", false
		}
		return parts[1], true
	}

	return "", false
}
