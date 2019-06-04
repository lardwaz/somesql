package somesql

import (
	"database/sql"
)

func exec(sql string, values []interface{}, db *sql.DB, autocommit bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	return execTx(sql, values, tx, autocommit)
}

func execTx(sql string, values []interface{}, tx *sql.Tx, autocommit bool) error {
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
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
