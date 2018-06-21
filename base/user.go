package base

import "time"

type User struct {
	id            string
	utype         int
	nickname      string
	password      string
	lastHeartBeat time.Time
}

var hbGapOut = time.Duration(time.Second * 3)

func NewUser(id string) *User {
	u := new(User)
	u.id = id
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
	return u.id
}

func (u *User) NickName() string {
	return u.nickname
}

func (u *User) Type() int {
	return u.utype
}
