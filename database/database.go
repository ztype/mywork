package database

import (
	"fmt"
	"io/ioutil"
	"log"

	"mywork/base"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const dbPath = `redten@mysql`

type DB struct {
	db *gorm.DB
}

var UserNotFound = fmt.Errorf("user not found")

func readToken() string {
	token, err := ioutil.ReadFile("./db.conf")
	if err != nil {
		log.Println("db conf", err)
		return ""
	}
	t := string(token)
	log.Println(t)
	return t
}

func NewDatabase() (*DB, error) {
	sdb, err := gorm.Open("mysql", readToken())
	sdb.LogMode(false)
	db := &DB{sdb}
	return db, err
}

func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) createTableUser() error {
	u := new(base.User)
	return db.db.CreateTable(u).Error
}

func (db *DB) dropTableUser() error {
	return db.db.DropTable(new(base.User)).Error
	return nil
}

//func initDatabse(db *gorm.DB) error {
//	info, _ := os.Lstat(dbPath)
//	if info.Size() == 0 {
//		return createUserTable(db)
//	}
//	return nil
//}
//
func (db *DB) Init() error {
	return nil
}

func (db *DB) InsertUser(user *base.User) error {
	//	stmt, err := db.db.Exec(`INSERT INTO user
	//(uid,nickname,isonline,regtime)
	//VALUES
	//(?,?,?,?) `)
	//if err != nil {
	//	return err
	//}
	//ret, err := stmt.Exec(user.Id(), user.NickName(), user.IsOnline(), time.Now().Unix())
	return db.db.Create(user).Error
}

//func (db *DB) UpdateUserOnline(id string, isonline bool) error {
//	glock.Lock()
//	defer glock.Unlock()
//	stmt, err := db.db.Prepare(`UPDATE user SET isonlien=? WHERE uid=?;`)
//	if err != nil {
//		return err
//	}
//	ret, err := stmt.Exec(isonline, id)
//	_ = ret
//	return err
//}
//
//func (db *DB) GetUserById(id string) (*base.User, error) {
//	glock.Lock()
//	glock.Unlock()
//	rows, err := db.db.Query(`SELECT
//uid,
//IFNULL(nickname,""),
//IFNULL(utype,0),
//IFNULL(headurl,"") FROM user WHERE uid=?;`, id)
//	if err != nil {
//		return nil, err
//	}
//	user := new(base.User)
//	for rows.Next() {
//		if err := rows.Scan(&user.Uid, &user.Nickname, &user.Utype, &user.Headurl); err != nil {
//			return nil, err
//		}
//		return user, nil
//	}
//	return nil, UserNotFound
//}
