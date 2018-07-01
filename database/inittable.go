package database

import (
	"github.com/jinzhu/gorm"
)

const usertablesql = `CREATE TABLE user (uid TEXT PRIMARY KEY NOT NULL,
nickname TEXT,
password TEXT,
utype INTEGER default 0,
isonline integer default 0,
headurl TEXT,
regtime INTEGER );`

func createUserTable(db *gorm.DB) error {
	db.CreateTable()

	return nil
}
