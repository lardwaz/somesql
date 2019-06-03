package somesql

import (
	"database/sql"
	"fmt"
	"strings"
)

// Select generates Postgres SELECT statement
// Implements: Mutator
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
		dataFieldLang = GetFieldData(s.GetLang())
	)

	metaFields := make([]string, 0)
	dataFields := make([]string, 0)
	relFields := make([]string, 0)
	for _, f := range s.fields {
		if innerField := GetInnerDataField(f); innerField != "" {
			if isInnerQuery {
				dataFields = append(dataFields, fmt.Sprintf(`"%s"->>'%s' "%s"`, dataFieldLang, innerField, innerField))
			} else {
				dataFields = append(dataFields, fmt.Sprintf(`'%s', "%s"->'%s'`, innerField, dataFieldLang, innerField))
			}
		} else if innerField := GetInnerRelationsField(f); innerField != "" {
			if isInnerQuery {
				relFields = append(relFields, fmt.Sprintf(`"%s"->>'%s' "%s"`, FieldRelations, innerField, innerField))
			} else {
				relFields = append(relFields, fmt.Sprintf(`'%s', "%s"->'%s'`, innerField, FieldRelations, innerField))
			}
		} else if IsFieldMeta(f) {
			if f == FieldData {
				f = dataFieldLang
			}
			metaFields = append(metaFields, fmt.Sprintf(`"%s"`, f))
		}
	}

	fields := make([]string, 0)
	if len(metaFields) > 0 {
		fields = append(fields, strings.Join(metaFields, ", "))
	}

	if len(dataFields) > 0 {
		if isInnerQuery {
			fields = append(fields, strings.Join(dataFields, ", "))
		} else {
			fields = append(fields, fmt.Sprintf(`json_build_object(%s) "%s"`, strings.Join(dataFields, ", "), FieldData))
		}
	}

	if len(relFields) > 0 {
		if isInnerQuery {
			fields = append(fields, strings.Join(relFields, ", "))
		} else {
			fields = append(fields, fmt.Sprintf(`json_build_object(%s) "%s"`, strings.Join(relFields, ", "), FieldRelations))
		}
	}

	fieldsStr = strings.Join(fields, ", ")

	conditions, values := processConditions(s.conditions)
	s.values = values

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
		s.sql = processPlaceholders(sql)
	}

	s.sql = cleanStatement(s.sql)
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
