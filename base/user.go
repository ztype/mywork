package base

import (
	"time"
)

type User struct {
	Uid           string `gorm:"unique;not null;index:user_id"`
	Utype         int
	Nickname      string
	Password      string
	Headurl       string
	lastHeartBeat time.Time
}

var hbGapOut = time.Duration(time.Second * 3)

func NewUser(id string) *User {
	u := new(User)
	u.Uid = id
	return u
}

func (u *User) IsOnline() bool {
	if time.Now().Sub(u.lastHeartBeat) > hbGapOut {
		return false
	}
	return true
}

func (u *User) HeartBeat() {
	u.lastHeartBeat = time.Now()
}

func (u *User) Id() string {
	return u.Uid
}

func (u *User) NickName() string {
	return u.Nickname
}

func (u *User) Type() int {
	return u.Utype
}
