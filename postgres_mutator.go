package somesql

import (
	"database/sql"
)

func rows(sql string, values []interface{}, db *sql.DB, autocommit bool) (*sql.Rows, error) {
	// TODO: implement
	return nil, nil
}

func rowsTx(sql string, values []interface{}, tx *sql.Tx, autocommit bool) (*sql.Rows, error) {
	// TODO: implement
	return nil, nil
}
