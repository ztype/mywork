package database

import (
	"io/ioutil"
	"log"

	"mywork/base"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DB struct {
	db *gorm.DB
}

var UserNotFound = gorm.ErrRecordNotFound

func readToken() string {
	token, err := ioutil.ReadFile("./db.conf")
	if err != nil {
		log.Println("db conf", err)
		return ""
	}
	t := string(token)
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
	if db.db.HasTable(new(base.User)) {
		return nil
	}
	return db.db.CreateTable(new(base.User)).Error
}

func (db *DB) dropTableUser() error {
	return db.db.DropTable(new(base.User)).Error
}

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

func (db *DB) GetUserById(id string) (*base.User, error) {
	u := new(base.User)
	u.Uid = id
	err := db.db.First(u).Error
	return u, err
}
