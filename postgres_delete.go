package somesql

import (
	"database/sql"
	"strconv"
)

// Delete generates Postgres DELETE statement
// Implements: Mutator
type Delete struct {
	conditions []Condition
	offset     int
	limit      int
	sql        string
	values     []interface{}
	db         *sql.DB
	lang       string
}

// NewDelete returns a new Delete
func NewDelete(db ...*sql.DB) *Delete {
	var s Delete

	s.lang = DefaultLang

	if len(db) > 0 {
		s.db = db[0]
	}

	return &s
}

// NewDeleteLang returns a new Delete with specific lang
func NewDeleteLang(lang string, db ...*sql.DB) *Delete {
	s := NewDelete(db...)
	s.SetLang(lang)

	return s
}

// SetDB implements Statement
func (s *Delete) SetDB(db *sql.DB) {
	s.db = db
}

// GetDB implements Statement
func (s Delete) GetDB() *sql.DB {
	return s.db
}

// SetLang implements Statement
func (s *Delete) SetLang(lang string) {
	s.lang = lang
}

// GetLang implements Statement
func (s Delete) GetLang() string {
	return s.lang
}

// GetSQL implements Statement
func (s Delete) GetSQL() string {
	return s.sql
}

// GetValues implements Statement
func (s Delete) GetValues() []interface{} {
	return s.values
}

// ToSQL implements Statement
func (s *Delete) ToSQL() {
	var (
		conditionsStr string
		offsetStr     string
		limitStr      string
	)

	if s.GetLang() == "" {
		s.SetLang(DefaultLang)
	}

	conditions, values := processConditions(s.conditions)
	s.values = values

	if len(conditions) > 0 {
		conditionsStr = "WHERE " + conditions
	}

	if s.limit > 0 {
		limitStr = "LIMIT " + strconv.Itoa(s.limit)
	}

	if s.offset > 0 {
		offsetStr = "OFFSET " + strconv.Itoa(s.offset)
	}

	sql := "DELETE FROM " + Table + " " + conditionsStr + " " + limitStr + " " + offsetStr

	s.sql = cleanStatement(processPlaceholders(sql))
}

// Exec implements Mutator
func (s Delete) Exec(autocommit bool) error {
	if s.GetSQL() == "" || len(s.GetValues()) == 0 {
		s.ToSQL()
	}

	return exec(s.GetSQL(), s.GetValues(), s.GetDB(), autocommit)
}

// ExecTx implements Mutator
func (s Delete) ExecTx(tx *sql.Tx, autocommit bool) error {
	if s.GetSQL() == "" || len(s.GetValues()) == 0 {
		s.ToSQL()
	}

	return execTx(s.GetSQL(), s.GetValues(), tx, autocommit)
}

// Where adds a condition clause to the Query
func (s *Delete) Where(c Condition) *Delete {
	s.conditions = append(s.conditions, c)
	return s
}

// Offset sets the Offset for Delete
func (s *Delete) Offset(offset int) *Delete {
	s.offset = offset
	return s
}

// Limit sets the Limit for Delete
func (s *Delete) Limit(limit int) *Delete {
	s.limit = limit
	return s
}
