package database

import "database/sql"

const usertablesql = `CREATE TABLE user(uid TEXT PRIMARY KEY NOT NULL,
nickname TEXT,
password TEXT,
utype INTEGER default 0,
isonline integer default 0,
headurl TEXT,
regtime INTEGER );`

func createUserTable(db *sql.DB) error {
	stmt, err := db.Prepare(usertablesql)
	if err != nil {
		return err
	}
	ret, err := stmt.Exec()
	if err != nil {
		return err
	}
	_ = ret
	return nil
}
