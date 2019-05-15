package somesql

import (
	"database/sql"
	"errors"
)

// QueryResulter related errors
var (
	ErrorNewTxCreate = errors.New("error creating new Tx")
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
func (q QueryResult) Exec(db *sql.DB, autocommit bool) error {
	var (
		err error
	)

	tx := q.GetTx()
	if tx == nil && db == nil {
		return ErrorNewTxCreate
	} else if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	stmt, err := tx.Prepare(q.GetSQL())
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(q.Values...)
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
func (q QueryResult) Rows(db *sql.DB) (*sql.Rows, error) {
	var (
		err error
	)

	tx := q.GetTx()
	if tx == nil && db == nil {
		return nil, ErrorNewTxCreate
	} else if tx == nil {
		tx, err = db.Begin()
		if err != nil {
			return nil, err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	rows, err := tx.Query(q.GetSQL(), q.Values...)
	if err != nil {
		return nil, err
	}
	// defer rows.Close() // Is it safe? We closed (defer) then returned it

	return rows, nil
}
