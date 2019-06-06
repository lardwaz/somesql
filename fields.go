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

	NoneJSONBArr uint8 = iota
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
	act := NoneJSONBArr
	if len(action) == 1 {
		act = action[0]
	}

	// if key already exist append to previous value
	var newVal []interface{}
	newValSlice, wasNewValAsserted := expandValues(value)
	jsonbField, exist := j.data[field]
	if exist {
		oldValSlice, isOldValSlice := j.data[field].Value.([]interface{})
		newValSliceInterface, isNewValSliceInterface := value.([]interface{})

		if isOldValSlice && isNewValSliceInterface {
			newVal = append(oldValSlice, newValSliceInterface...)
		} else if isOldValSlice {
			newVal = append(oldValSlice, value)
		} else if isNewValSliceInterface {
			newVal = append(newVal, jsonbField.Value)
			newVal = append(newVal, newValSliceInterface...)
		} else {
			newVal = []interface{}{jsonbField.Value, value}
		}

		j.data[field] = JSONBField{Value: newVal, Action: act}
	} else if wasNewValAsserted {
		j.data[field] = JSONBField{Value: newValSlice, Action: act}
	} else {
		j.data[field] = JSONBField{Value: value, Action: act}
		j.keys = append(j.keys, field)
	}
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
	f.set(field, value, NoneJSONBArr)

	return f
}

// Add adds data to an array jsonb field
func (f Fields) Add(field string, value interface{}) Fields {
	f.set(field, value, JSONBArrAdd)

	return f
}

// Remove removes data from an array jsonb field
func (f Fields) Remove(field string, value interface{}) Fields {
	f.set(field, value, JSONBArrRemove)

	return f
}

func (f Fields) set(field string, value interface{}, action uint8) {
	if IsFieldMeta(field) || IsFieldData(field) || IsFieldRelations(field) {
		f[field] = value
	} else if innerField, ok := GetInnerField(FieldData, field); ok {
		jsonbFields, ok := f[FieldData].(JSONBFields)
		if !ok { // if not jsonbfields, make it
			jsonbFields = NewJSONBFields()
		}
		jsonbFields.Add(innerField, value, action)
		f[FieldData] = jsonbFields
	} else if innerField, ok := GetInnerField(FieldRelations, field); ok {
		jsonbFields, ok := f[FieldRelations].(JSONBFields)
		if !ok { // if not jsonbfields, make it
			jsonbFields = NewJSONBFields()
		}
		vals, _ := expandValues(value)
		jsonbFields.Add(innerField, vals, action)
		f[FieldRelations] = jsonbFields
	}
}

// List returns fields and values
func (f Fields) List() ([]string, []interface{}) {
	var (
		fields []string
		values []interface{}
	)

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
	switch field {
	case FieldID, FieldCreatedAt, FieldUpdatedAt, FieldOwnerID, FieldStatus, FieldType:
		return true
	}

	return false
}

// IsFieldData returns true if field is a data field
func IsFieldData(field string) bool {
	return field == FieldData
}

// IsFieldRelations returns true if field is a data field
func IsFieldRelations(field string) bool {
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
