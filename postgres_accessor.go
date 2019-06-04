package somesql

import (
	"database/sql"
)

func rows(sql string, values []interface{}, db *sql.DB) (*sql.Rows, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	return rowsTx(sql, values, tx)
}

func rowsTx(sql string, values []interface{}, tx *sql.Tx) (*sql.Rows, error) {
	rows, err := tx.Query(sql, values...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
