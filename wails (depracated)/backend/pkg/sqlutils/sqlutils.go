package sqlutils

import "github.com/jmoiron/sqlx"

func ExecStruct[T comparable](tx *sqlx.Tx, sqlQuery string, data T) error {
	smt, err := tx.PrepareNamed(sqlQuery)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(data)
	return err
}
