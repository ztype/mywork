package database

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"database/sql"
	"log"
)

func Test_open(t *testing.T) {
	db, err := sql.Open("sqlite3", "E:/t.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	t.Log(db.Ping())
	err = createTable(db)
	t.Log(err)
}

func createTable(db *sql.DB) error {
	s := `CREATE TABLE login(uid TEXT PRIMARY KEY NOT NULL,lasttime INTEGER,address TEXT);`
	stmt, err := db.Prepare(s)
	if err != nil {
		return err
	}
	ret, err := stmt.Exec()
	if err != nil {
		return err
	}
	log.Println(ret)
	return nil
}
