package somesql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// Insert generates Postgres INSERT statement
// Implements: Mutator
type Insert struct {
	fields Fields
	sql    string
	values []interface{}
	db     *sql.DB
	lang   string
}

// NewInsert returns a new Insert
func NewInsert(lang string, db ...*sql.DB) *Insert {
	var s Insert

	s.lang = lang
	s.fields = NewFields()

	if len(db) > 0 {
		s.db = db[0]
	}

	return &s
}

// SetDB implements Statement
func (s *Insert) SetDB(db *sql.DB) {
	s.db = db
}

// GetDB implements Statement
func (s Insert) GetDB() *sql.DB {
	return s.db
}

// SetLang implements Statement
func (s *Insert) SetLang(lang string) {
	s.lang = lang
}

// GetLang implements Statement
func (s Insert) GetLang() string {
	return s.lang
}

// GetSQL implements Statement
func (s Insert) GetSQL() string {
	return s.sql
}

// GetValues implements Statement
func (s Insert) GetValues() []interface{} {
	return s.values
}

// ToSQL implements Statement
func (s *Insert) ToSQL() {
	var (
		placeholderIndex int
		dataFieldLang    = GetFieldData(s.GetLang())
	)
	fields, values := s.fields.List()

	// Processing fields and values
	metaFields := make([]string, 0)
	dataFields := make([]string, 0)
	relFields := make([]string, 0)
	metaValues := make([]interface{}, 0)
	dataValues := make(map[string]interface{})
	relValues := make(map[string]interface{})
	for i, f := range fields {
		if innerField := GetInnerDataField(f); innerField != "" { // Check if Inner Data + Relations fields
			dataFields = append(dataFields, fmt.Sprintf(`"%s"`, innerField))
			dataValues[innerField] = values[i]
		} else if innerField := GetInnerRelationsField(f); innerField != "" {
			relFields = append(relFields, fmt.Sprintf(`"%s"`, innerField))
			relValues[innerField] = values[i]
		} else if IsFieldData(f) { // Check if Outer Data + Relations fields
			placeholderIndex++
			metaFields = append(metaFields, fmt.Sprintf(`"%s"`, dataFieldLang))
			if vals, err := json.Marshal(values[i]); err == nil {
				metaValues = append(metaValues, string(vals))
			}
		} else if IsFieldRelations(f) {
			placeholderIndex++
			metaFields = append(metaFields, fmt.Sprintf(`"%s"`, f))
			if vals, err := json.Marshal(values[i]); err == nil {
				metaValues = append(metaValues, string(vals))
			}
		} else if IsFieldMeta(f) { // Check if Meta fields
			placeholderIndex++
			metaFields = append(metaFields, fmt.Sprintf(`"%s"`, f))
			metaValues = append(metaValues, values[i])
		}
	}

	fieldsJoined := make([]string, 0)
	if len(metaFields) > 0 {
		fieldsJoined = append(fieldsJoined, strings.Join(metaFields, ", "))
		s.values = append(s.values, metaValues...)
	}

	if len(dataFields) > 0 {
		placeholderIndex++
		fieldsJoined = append(fieldsJoined, fmt.Sprintf(`"%s"`, dataFieldLang))
		if vals, err := json.Marshal(dataValues); err == nil {
			s.values = append(s.values, string(vals))
		}
	}

	if len(relFields) > 0 {
		placeholderIndex++
		fieldsJoined = append(fieldsJoined, fmt.Sprintf(`"%s"`, FieldRelations))
		if vals, err := json.Marshal(relValues); err == nil {
			s.values = append(s.values, string(vals))
		}
	}

	fieldsStr := strings.Join(fieldsJoined, ", ")

	placeholders := make([]string, 0)
	for i := 1; i <= placeholderIndex; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
	}
	placesholdersStr := strings.Join(placeholders, ", ")

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, Table, fieldsStr, placesholdersStr)

	s.sql = cleanStatement(sql)
}

// Exec implements Mutator
func (s Insert) Exec(autocommit bool) error {
	return exec(s.GetSQL(), s.GetValues(), s.GetDB(), autocommit)
}

// ExecTx implements Mutator
func (s Insert) ExecTx(tx *sql.Tx, autocommit bool) error {
	return execTx(s.GetSQL(), s.GetValues(), tx, autocommit)
}

// Fields sets the fields and values for insert
func (s *Insert) Fields(fields Fields) *Insert {
	s.fields = fields
	return s
}
