package somesql

import (
	"database/sql"
	"encoding/json"
	"strconv"
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

	s.fields = NewFields()
	s.lang = lang

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
		fieldsStr        string
		placesholdersStr string
		placeholderIndex int
		dataFieldLang    string

		fieldsBuff       strings.Builder
		placeholdersBuff strings.Builder
	)

	dataFieldLang = GetLangFieldData(s.GetLang())

	fields, values := s.fields.List()

	// Processing fields and values
	s.values = make([]interface{}, len(fields))
	for i, f := range fields {
		if IsFieldMeta(f) {
			s.values[i] = values[i]
		} else if IsFieldData(f) || IsFieldRelations(f) {
			if jsonbFields, ok := values[i].(JSONBFields); ok {
				if jsonBytes, err := json.Marshal(jsonbFields.Values()); err == nil {
					s.values[i] = string(jsonBytes)
				}
			}
		}

		// Double quote the field name
		// Placeholders
		if IsFieldMeta(f) || IsFieldData(f) || IsFieldRelations(f) {
			if IsFieldData(f) || IsFieldRelations(f) {
				f = dataFieldLang // data => data_<lang>
			}
			fieldsBuff.WriteString(`"` + f + `", `)
			placeholderIndex++
			placeholdersBuff.WriteString(`$` + strconv.Itoa(placeholderIndex) + `, `)
		}
	}

	if fieldsBuff.Len() > 0 {
		fieldsStr = fieldsBuff.String()[:fieldsBuff.Len()-2] // trim ", "
	}

	if placeholdersBuff.Len() > 0 {
		placesholdersStr = placeholdersBuff.String()[:placeholdersBuff.Len()-2] // trim ", "
	}

	sql := "INSERT INTO " + Table + " (" + fieldsStr + ") VALUES (" + placesholdersStr + ")"

	s.sql = cleanStatement(sql)
}

// Exec implements Mutator
func (s Insert) Exec(autocommit bool) error {
	if s.GetSQL() == "" || len(s.GetValues()) == 0 {
		s.ToSQL()
	}

	return exec(s.GetSQL(), s.GetValues(), s.GetDB(), autocommit)
}

// ExecTx implements Mutator
func (s Insert) ExecTx(tx *sql.Tx, autocommit bool) error {
	if s.GetSQL() == "" || len(s.GetValues()) == 0 {
		s.ToSQL()
	}

	return execTx(s.GetSQL(), s.GetValues(), tx, autocommit)
}

// Fields sets the fields and values for insert
func (s *Insert) Fields(fields Fields) *Insert {
	s.fields = fields
	return s
}
