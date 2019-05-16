package somesql

import (
	"database/sql"
	"errors"
)

// QueryResulter related errors
var (
	ErrorNoDBTX = errors.New("invalid DB or Tx")
)

// QueryResult is an implementation of QueryResulter
type QueryResult struct {
	Query
	SQL    string
	Values []interface{}
}

// NewQueryResult returns a new QueryResult
func NewQueryResult(query Query, sql string, values []interface{}) QueryResult {
	return QueryResult{Query: query, SQL: sql, Values: values}
}

// GetSQL returns the sql stmt
func (q QueryResult) GetSQL() string {
	return q.SQL
}

// GetValues returns the list of values
func (q QueryResult) GetValues() []interface{} {
	return q.Values
}

// Exec executes stmt values using a specific sql transaction
func (q QueryResult) Exec(autocommit bool) error {
	db := q.GetDB()
	if db == nil {
		return ErrorNoDBTX
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	return q.ExecTx(tx, autocommit)
}

// ExecTx executes stmt values using a specific sql transaction
func (q QueryResult) ExecTx(tx *sql.Tx, autocommit bool) error {
	stmt, err := tx.Prepare(q.GetSQL())
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(q.GetValues()...)
	if err != nil {
		return err
	}

	if autocommit {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// Rows executes stmt values using a specific sql transaction
func (q QueryResult) Rows() (*sql.Rows, error) {
	db := q.GetDB()
	if db == nil {
		return nil, ErrorNoDBTX
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	return q.RowsTx(tx)
}

// RowsTx executes stmt values using a specific sql transaction
func (q QueryResult) RowsTx(tx *sql.Tx) (*sql.Rows, error) {
	rows, err := tx.Query(q.GetSQL(), q.GetValues()...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
