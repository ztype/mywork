package database

import (
	"database/sql"
	"mywork/base"

	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = `./redten.db`

type DB struct {
	*sql.DB
}

func NewDatabase() (*DB, error) {
	sdb, err := sql.Open("sqlite3", dbPath)
	db := &DB{sdb}
	return db, err
}

func (db *DB) Init() error {
	return nil
}

func (db *DB) InsertUser(user *base.User) error {
	stmt, err := db.Prepare(`INSERT INTO user (uid,nickname,isonline,regtime) VALUES (?,?,?,?) `)
	if err != nil {
		return err
	}
	ret, err := stmt.Exec(user.Id(), user.NickName(), user.IsOnline(), time.Now().Unix())
	_ = ret
	return err
}

func (db *DB) UpdateUserOnline(user *base.User) error {
	stmt, err := db.Prepare(`UPDATE user SET isonlien=? WHERE uid=?;`)
	if err != nil {
		return err
	}
	ret, err := stmt.Exec(user.IsOnline(), user.Id())
	_ = ret
	return err
}

func (db *DB) GetUserById(id string) (*base.User, error) {
	stmt, err := db.Query(`SELECT * FROM user WHERE uid=?;`)
	if err != nil {
		return nil, err
	}

}
