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

func SetHeartbeatTime(d time.Duration) {
	hbGapOut = d
}

func NewUser(id string) *User {
	u := new(User)
	u.Uid = id
	return u
}

func (u *User) IsOnline() bool {
	return time.Now().Sub(u.lastHeartBeat) < hbGapOut

}

func (u *User) HeartBeat() {
	u.lastHeartBeat = time.Now()
}

func (u *User) ID() string {
	return u.Uid
}

func (u *User) NickName() string {
	return u.Nickname
}

func (u *User) Type() int {
	return u.Utype
}
