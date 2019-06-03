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
	var dataFieldLang = GetFieldData(s.GetLang())
	fields, values := s.fields.List()

	// Processing fields and values
	placeholders := make([]string, len(fields))
	for i, f := range fields {
		if IsFieldData(f) || IsFieldRelations(f) {
			if IsFieldData(f) {
				f = dataFieldLang
			}
			if json, err := json.Marshal(values[i]); err == nil {
				values[i] = string(json)
			}
		}
		fields[i] = fmt.Sprintf(`"%s"`, f)
		placeholders[i] = fmt.Sprintf(`$%d`, i+1)
	}

	fieldsStr := strings.Join(fields, ", ")
	placesholdersStr := strings.Join(placeholders, ", ")

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, Table, fieldsStr, placesholdersStr)

	s.sql = cleanStatement(sql)
	s.values = values
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
