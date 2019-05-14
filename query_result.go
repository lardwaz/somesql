package somesql

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
	tx := q.GetTx()
	stmt, err := tx.Prepare(q.SQL)
	defer stmt.Close()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(q.Values...)

	return err
}
