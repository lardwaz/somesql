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
	// fields, values := s.fields.List()

	// // Processing fields and values
	// placeholders := make([]string, len(fields))
	// for i, f := range fields {
	// 	if IsFieldData(f) || IsFieldRelations(f) {
	// 		if IsFieldData(f) {
	// 			f = GetFieldData(s.GetLang())
	// 		}
	// 		if json, err := json.Marshal(values[i]); err == nil {
	// 			values[i] = string(json)
	// 		}
	// 	}
	// 	fields[i] = fmt.Sprintf(`"%s"`, f)
	// 	placeholders[i] = fmt.Sprintf(`$%d`, i+1)
	// }

	var (
		conditionsStr string
		offsetStr     string
		limitStr      string
	)

	// TODO: Processing fields and values

	var conditions string
	for i, cond := range s.conditions {
		if i != 0 {
			switch cond.ConditionType() {
			case AndCondition:
				conditions += ` AND `
			case OrCondition:
				conditions += ` OR `
			default:
				continue
			}
		}

		c, v := cond.AsSQL()
		conditions += c
		s.values = append(s.values, v...)
	}

	fieldsStr := strings.Join(s.fields, ", ")

	if len(conditions) > 0 {
		conditionsStr = fmt.Sprintf("WHERE %s", conditions)
	}

	if s.offset > 0 {
		offsetStr = fmt.Sprintf("OFFSET %d", s.offset)
	}

	if s.limit > 0 {
		limitStr = fmt.Sprintf("LIMIT %d", s.limit)
	}

	sql := fmt.Sprintf(`SELECT %s FROM %s %s %s %s`, fieldsStr, Table, conditionsStr, offsetStr, limitStr)

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
	s.fields = fields
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
