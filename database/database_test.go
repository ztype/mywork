package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_open(t *testing.T) {
	db, err := NewDatabase()
	if err != nil {
		t.Fatal(err)
	}
	//err = db.createTableUser()
	err = db.dropTableUser()
	t.Log(err)
	db.Close()
}
