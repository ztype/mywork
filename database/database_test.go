package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_open(t *testing.T) {
	db, err := sql.Open("sqlite3", "E:/t.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	t.Log(db.Ping())
	err = createUserTable(db)
	t.Log(err)
}
