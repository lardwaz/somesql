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
		dataFieldLang    = GetLangFieldData(s.GetLang())
	)
	fields, values := s.fields.List()

	// Processing fields and values
	s.values = make([]interface{}, len(fields))
	for i, f := range fields {
		if IsWholeFieldData(f) {
			if jsonbFields, ok := values[i].(JSONBFields); ok {
				if jsonBytes, err := json.Marshal(jsonbFields.Values()); err == nil {
					s.values[i] = string(jsonBytes)
				}
				f = dataFieldLang // data => data_<lang>
				placeholderIndex++
			}
		} else if IsWholeFieldRelations(f) {
			if jsonbFields, ok := values[i].(JSONBFields); ok {
				if jsonBytes, err := json.Marshal(jsonbFields.Values()); err == nil {
					s.values[i] = string(jsonBytes)
				}
				placeholderIndex++
			}
		} else if IsFieldMeta(f) { // Check if Meta fields
			s.values[i] = values[i]
			placeholderIndex++
		}

		// Double quote the field name
		fields[i] = fmt.Sprintf(`"%s"`, f)
	}

	fieldsStr := strings.Join(fields, ", ")

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
