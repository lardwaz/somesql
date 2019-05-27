package somesql

import (
	"database/sql"
)

func exec(sql string, values []interface{}, db *sql.DB, autocommit bool) error {
	// TODO: implement
	return nil
}

func execTx(sql string, values []interface{}, tx *sql.Tx, autocommit bool) error {
	// TODO: implement
	return nil
}
