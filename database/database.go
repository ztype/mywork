package database

import (
	"database/sql"
	"mywork/base"
	"time"

	"fmt"

	"os"

	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var glock sync.Mutex

const dbPath = `./redten.db`

type DB struct {
	db *sql.DB
}

var UserNotFound = fmt.Errorf("user not found")

func NewDatabase() (*DB, error) {
	sdb, err := sql.Open("sqlite3", dbPath)
	db := &DB{sdb}
	if err := initDatabse(sdb); err != nil {
		return nil, err
	}
	return db, err
}

func initDatabse(db *sql.DB) error {
	info, _ := os.Lstat(dbPath)
	if info.Size() == 0 {
		return createUserTable(db)
	}
	return nil
}

func (db *DB) Init() error {
	return nil
}

func (db *DB) InsertUser(user *base.User) error {
	glock.Lock()
	defer glock.Unlock()
	stmt, err := db.db.Prepare(`INSERT INTO user 
(uid,nickname,isonline,regtime) 
VALUES 
(?,?,?,?) `)
	if err != nil {
		return err
	}
	ret, err := stmt.Exec(user.Id(), user.NickName(), user.IsOnline(), time.Now().Unix())
	_ = ret
	return err
}

func (db *DB) UpdateUserOnline(id string, isonline bool) error {
	glock.Lock()
	defer glock.Unlock()
	stmt, err := db.db.Prepare(`UPDATE user SET isonlien=? WHERE uid=?;`)
	if err != nil {
		return err
	}
	ret, err := stmt.Exec(isonline, id)
	_ = ret
	return err
}

func (db *DB) GetUserById(id string) (*base.User, error) {
	glock.Lock()
	glock.Unlock()
	rows, err := db.db.Query(`SELECT 
uid,
IFNULL(nickname,""),
IFNULL(utype,0),
IFNULL(headurl,"") FROM user WHERE uid=?;`, id)
	if err != nil {
		return nil, err
	}
	user := new(base.User)
	for rows.Next() {
		if err := rows.Scan(&user.Uid, &user.Nickname, &user.Utype, &user.Headurl); err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, UserNotFound
}
