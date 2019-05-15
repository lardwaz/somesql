package somesql

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// FieldValue implements the FieldValuer interface
type FieldValue struct {
	pos       map[string]int
	fields    []string
	values    []interface{}
	relations map[string][]string
}

// NewFieldValue returns a new FieldValue
func NewFieldValue() *FieldValue {
	var f FieldValue

	f.pos = make(map[string]int, 0)
	f.fields = make([]string, 0)
	f.values = make([]interface{}, 0)
	f.relations = make(map[string][]string, 0)

	return &f
}

// ID is a setter for ID in FieldValue
func (f *FieldValue) ID(id string) FieldValuer {
	f.Set("id", id)

	return f
}

// CreatedAt is a setter for CreatedAt in FieldValue
func (f *FieldValue) CreatedAt(t time.Time) FieldValuer {
	f.Set("created_at", t)

	return f
}

// UpdatedAt is a setter for UpdatedAt in FieldValue
func (f *FieldValue) UpdatedAt(t time.Time) FieldValuer {
	f.Set("updated_at", t)

	return f
}

// OwnerID is a setter for OwnerID in FieldValue
func (f *FieldValue) OwnerID(id string) FieldValuer {
	f.Set("owner_id", id)

	return f
}

// Status is a setter for Status in FieldValue
func (f *FieldValue) Status(s string) FieldValuer {
	f.Set("status", s)

	return f
}

// Type is a setter for Type in FieldValue
func (f *FieldValue) Type(s string) FieldValuer {
	f.Set("type", s)

	return f
}

// Data is a setter for Data in FieldValue
func (f *FieldValue) Data(json string) FieldValuer {
	f.Set("data", json)

	return f
}

// UseDefaults is a setter for UseDefaults in FieldValue
func (f *FieldValue) UseDefaults() FieldValuer {
	return f.ID(uuid.NewV4().String()).CreatedAt(time.Now()).UpdatedAt(time.Now()).OwnerID(uuid.Nil.String()).Status("").Type("").Data("{}")
}

// Set implements the FieldValuer interface
func (f *FieldValue) Set(field string, value interface{}) FieldValuer {
	if pos, ok := f.pos[field]; ok {
		f.values[pos] = value
		return f
	}

	f.fields = append(f.fields, field)
	f.values = append(f.values, value)
	f.pos[field] = len(f.fields) - 1

	return f
}

// SetRel implements the FieldValuer interface
func (f *FieldValue) SetRel(rel string, value []string) FieldValuer {
	f.relations[rel] = value
	return f
}

// List implements the FieldValuer interface
func (f *FieldValue) List() ([]string, []interface{}, map[string][]string) {
	return f.fields, f.values, f.relations
}
