package somesql

import (
	"database/sql"
	"errors"
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

	err = execTx(sql, values, tx, autocommit)

	return err
}

func execTx(sql string, values []interface{}, tx *sql.Tx, autocommit bool) error {
	if sql == "" || len(values) == 0 {
		return errors.New("invalid sql or values")
	}

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
