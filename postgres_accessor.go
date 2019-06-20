package somesql

import (
	"database/sql"
	"errors"
)

func rows(sql string, values []interface{}, db *sql.DB) (*sql.Rows, error) {
	if sql == "" || len(values) == 0 {
		return nil, errors.New("invalid sql or values")
	}

	rows, err := db.Query(sql, values...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
