package somesql

import (
	"database/sql"
	"fmt"
	"strings"
)

// Select generates Postgres SELECT statement
// Implements: Accessor
type Select struct {
	fields     []string
	conditions []Condition
	inner      bool
	offset     int
	limit      int
	sql        string
	values     []interface{}
	db         *sql.DB
	lang       string
}

// NewSelect returns a new Select
func NewSelect(lang string, inner bool, db ...*sql.DB) *Select {
	var s Select

	s.lang = lang
	s.inner = inner
	s.fields = FieldsList
	s.limit = 10

	if len(db) > 0 {
		s.db = db[0]
	}

	return &s
}

// SetDB implements Statement
func (s *Select) SetDB(db *sql.DB) {
	s.db = db
}

// GetDB implements Statement
func (s Select) GetDB() *sql.DB {
	return s.db
}

// SetLang implements Statement
func (s *Select) SetLang(lang string) {
	s.lang = lang
}

// GetLang implements Statement
func (s Select) GetLang() string {
	return s.lang
}

// GetSQL implements Statement
func (s Select) GetSQL() string {
	return s.sql
}

// GetValues implements Statement
func (s Select) GetValues() []interface{} {
	return s.values
}

// ToSQL implements Statement
func (s *Select) ToSQL() {
	var (
		fieldsStr     string
		conditionsStr string
		offsetStr     string
		limitStr      string
		isInnerQuery  = s.IsInner()
		dataFieldLang = GetLangFieldData(s.GetLang())
	)

	metaFields := make([]string, 0)
	dataFields := make([]string, 0)
	relFields := make([]string, 0)
	for _, f := range s.fields {
		if innerField, ok := GetInnerField(FieldData, f); ok {
			if isInnerQuery {
				dataFields = append(dataFields, fmt.Sprintf(`"%s"->>'%s' "%s"`, dataFieldLang, innerField, innerField))
			} else {
				dataFields = append(dataFields, fmt.Sprintf(`'%s', "%s"->'%s'`, innerField, dataFieldLang, innerField))
			}
		} else if innerField, ok := GetInnerField(FieldRelations, f); ok {
			if isInnerQuery {
				relFields = append(relFields, fmt.Sprintf(`"%s"->>'%s' "%s"`, FieldRelations, innerField, innerField))
			} else {
				relFields = append(relFields, fmt.Sprintf(`'%s', "%s"->'%s'`, innerField, FieldRelations, innerField))
			}
		} else if IsFieldMeta(f) || IsFieldData(f) || IsFieldRelations(f) {
			if f == FieldData {
				f = dataFieldLang
			}
			metaFields = append(metaFields, fmt.Sprintf(`"%s"`, f))
		}
	}

	fieldsJoined := make([]string, 0)
	// Meta fields
	if len(metaFields) > 0 {
		fieldsJoined = append(fieldsJoined, strings.Join(metaFields, ", "))
	}

	// Data fields
	if len(dataFields) > 0 {
		dataFieldsJoined := strings.Join(dataFields, ", ")
		if isInnerQuery {
			fieldsJoined = append(fieldsJoined, dataFieldsJoined)
		} else {
			fieldsJoined = append(fieldsJoined, fmt.Sprintf(`json_build_object(%s) "%s"`, dataFieldsJoined, FieldData))
		}
	}

	// Relationship fields
	if len(relFields) > 0 {
		relFieldsJoined := strings.Join(relFields, ", ")
		if isInnerQuery {
			fieldsJoined = append(fieldsJoined, relFieldsJoined)
		} else {
			fieldsJoined = append(fieldsJoined, fmt.Sprintf(`json_build_object(%s) "%s"`, relFieldsJoined, FieldRelations))
		}
	}

	fieldsStr = strings.Join(fieldsJoined, ", ")

	conditions, condValues := processConditions(s.conditions)
	s.values = condValues

	if len(conditions) > 0 {
		conditionsStr = fmt.Sprintf("WHERE %s", conditions)
	}

	if s.limit > 0 {
		limitStr = fmt.Sprintf("LIMIT %d", s.limit)
	}

	if s.offset > 0 {
		offsetStr = fmt.Sprintf("OFFSET %d", s.offset)
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s %s %s %s`, fieldsStr, Table, conditionsStr, limitStr, offsetStr)

	if !isInnerQuery {
		sql = processPlaceholders(sql)
	}

	s.sql = cleanStatement(sql)
}

// SetInner implements Accessor
func (s *Select) SetInner(inner bool) {
	s.inner = inner
}

// IsInner implements Accessor
func (s Select) IsInner() bool {
	return s.inner
}

// Rows implements Accessor
func (s Select) Rows() (*sql.Rows, error) {
	return rows(s.GetSQL(), s.GetValues(), s.GetDB())
}

// RowsTx implements Accessor
func (s Select) RowsTx(tx *sql.Tx) (*sql.Rows, error) {
	return rowsTx(s.GetSQL(), s.GetValues(), tx)
}

// Fields sets the fields for Select
func (s *Select) Fields(fields ...string) *Select {
	if len(fields) == 0 {
		return s
	}
	s.fields = fields
	return s
}

// Where adds a condition clause to the Query
func (s *Select) Where(c Condition) *Select {
	s.conditions = append(s.conditions, c)
	return s
}

// Offset sets the Offset for Select
func (s *Select) Offset(offset int) *Select {
	s.offset = offset
	return s
}

// Limit sets the Limit for Select
func (s *Select) Limit(limit int) *Select {
	s.limit = limit
	return s
}
