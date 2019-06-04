package somesql //import go.lsl.digital/gocipe/somesql

import (
	"database/sql"
)

const (
	// AndCondition represents a condition added to the query via AND keyword
	AndCondition uint8 = iota
	// OrCondition represents a condition added to the query via OR keyword
	OrCondition

	// LangEN represents the english language
	LangEN string = "en"
	// LangFR represents the french language
	LangFR string = "fr"

	// Table represnts the table name
	Table = "repo"
)

const (
	// None represents a simple way of explicitly specifying no value
	None = ""
)

// Statement represents a composable statement
// Can be consumed by Mutator or Accessor
type Statement interface {
	SetDB(db *sql.DB)
	GetDB() *sql.DB
	SetLang(lang string)
	GetLang() string
	GetSQL() string
	GetValues() []interface{}
	ToSQL()
}

// Mutator is any statement which modifies values in store
type Mutator interface {
	Statement
	Exec(autocommit bool) error
	ExecTx(tx *sql.Tx, autocommit bool) error
}

// Accessor is any statement which retrieves values from store
type Accessor interface {
	Statement
	SetInner(inner bool)
	IsInner() bool
	Rows() (*sql.Rows, error)
	RowsTx(tx *sql.Tx) (*sql.Rows, error)
}

// Condition represents a conditional clause in a statement
type Condition interface {
	ConditionType() uint8
	AsSQL(in ...bool) (string, []interface{})
}
