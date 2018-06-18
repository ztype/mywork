package services

import "time"

type User struct {
	id            string
	lastHeartBeat time.Time
}

var hbGapOut = time.Duration(time.Second * 3)

func (u *User) IsOnline() bool {
	if time.Now().Sub(u.lastHeartBeat) > hbGapOut {
		return false
	}
	return true
}

func (u *User) HeartBeat() {
	u.lastHeartBeat = time.Now()
}
