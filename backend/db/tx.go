package db

import "database/sql"

type Tx struct {
	*sql.Tx
	commited bool
}

func (t *Tx) Commit() error {
	err := t.Tx.Commit()
	if err != nil {
		return err
	}
	t.commited = true
	return nil
}

func (t *Tx) AutoRollback() error {
	if !t.commited {
		return t.Tx.Rollback()
	}
	return nil
}
