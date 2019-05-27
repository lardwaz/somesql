package somesql //import go.lsl.digital/gocipe/somesql

import (
	"database/sql"
	"time"
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

// Query represents a composable query
// Setters:  Select(), Save(), Delete(), Where(), SetLang(), SetLimit(), SetOffset(), SetRel()
// Getters:  GetLang(), GetLimit(), GetOffset(), GetTx(), IsInner(), AsSQL()
// Mutators: AddRel(), RemoveRel()
type Query interface {
	Insert(FieldValuer) Query
	Select(fields ...string) Query
	SelectRel(rels ...string) Query
	Update(FieldValuer) Query
	Delete() Query
	Where(Condition) Query
	SetLang(lang string) Query
	GetLang() string
	SetLimit(limit int) Query
	GetLimit() int
	SetOffset(offset int) Query
	GetOffset() int
	SetDB(db *sql.DB) Query
	GetDB() *sql.DB
	SetInner(inner bool) Query
	IsInner() bool
	AsSQL() QueryResulter
	InsertRel(rel string, values []string) Query
	AddRel(rel string, values []string) Query
	RemoveRel(rel string, values []string) Query
}

// QueryResulter is the result of running AsSQL on Query
type QueryResulter interface {
	Query
	GetSQL() string
	GetValues() []interface{}
	Exec(autocommit bool) error
	ExecTx(tx *sql.Tx, autocommit bool) error
	Rows() (*sql.Rows, error)
	RowsTx(tx *sql.Tx) (*sql.Rows, error)
}

// FieldValuer assigns a value to a field
type FieldValuer interface {
	ID(id string) FieldValuer
	CreatedAt(t time.Time) FieldValuer
	UpdatedAt(t time.Time) FieldValuer
	OwnerID(id string) FieldValuer
	Status(s string) FieldValuer
	Type(s string) FieldValuer
	Data(json string) FieldValuer
	UseDefaults() FieldValuer
	Set(field string, value interface{}) FieldValuer
	List() ([]string, []interface{} /*, map[string][]string*/)
}

// Condition represents a conditional clause in a statement
type Condition interface {
	ConditionType() uint8
	AsSQL(in ...bool) (string, []interface{})
}
