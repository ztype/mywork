package main

import (
	"sync"
	"time"
	"log"
)

type SessionManager struct {
	users map[string]*User
	lock  sync.Mutex
}

type User struct {
	id            string
	lastHeartBeat time.Time
}

type Session struct {
	user *User
}

var hbGap = time.Duration(time.Second * 60)
var hbGapOut = hbGap + time.Duration(time.Second*5)

func (u *User) IsOnline() bool {
	if time.Now().Sub(u.lastHeartBeat) > hbGapOut {
		return false
	}
	return true
}

func (u *User)HeartBeat(){
	u.lastHeartBeat = time.Now()
}

func NewManager() *SessionManager {
	m := new(SessionManager)
	m.users = make(map[string]*User, 0)
	return m
}

func (sm *SessionManager) UserConnect(uuid string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	_, ok := sm.users[uuid]
	if ok {
		return
	}
	user := new(User)
	user.lastHeartBeat = time.Now()
	user.id = uuid
	sm.users[uuid] = user
	time.AfterFunc(hbGap, func() {
		checkHeartbeat(user,sm)
	})
	log.Println(uuid,"connected")
}

func (sm *SessionManager) UserDisConnect(uuid string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if _, ok := sm.users[uuid]; ok {
		delete(sm.users, uuid)
		log.Println(uuid,"disconnected")
	}
}

func checkHeartbeat(user *User, sm *SessionManager) {
	if !user.IsOnline() {
		sm.UserDisConnect(user.id)
	}else {
		time.AfterFunc(hbGap,func(){
			checkHeartbeat(user,sm)
		})
	}
}

func (sm *SessionManager) GetUser(uuid string) *User {
	if u, ok := sm.users[uuid]; ok {
		return u
	}
	return nil
}

func (sm *SessionManager) IsExist(uuid string) bool {
	_, ok := sm.users[uuid]
	return ok
}

func (sm *SessionManager) UserCount() int {
	return len(sm.users)
}
