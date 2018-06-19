package services

import (
	"log"
	"mywork/utils"
	"sync"
	"time"
)

const (
	heartbeat = "heartbeat"
	logout    = "logout"
)

type SessionManager struct {
	users map[string]*User
	lock  sync.Mutex
}

type Session struct {
	user *User
}

func NewManager() *SessionManager {
	m := new(SessionManager)
	m.users = make(map[string]*User, 0)
	go m.check()
	return m
}

func (sm *SessionManager) Name() string {
	return "connect"
}

func (sm *SessionManager) Serve(p utils.Param) (interface{}, error) {
	switch p.Type {
	case heartbeat:
		return sm.UserConnect(p.Id)
	case logout:
		return sm.UserDisConnect(p.Id)
	}
	return nil, nil
}

func (sm *SessionManager) check() {
	for {
		sm.checkUser()
		time.Sleep(hbGapOut)
	}
}

func (sm *SessionManager) checkUser() {
	for id, user := range sm.users {
		if !user.IsOnline() {
			sm.lock.Lock()
			delete(sm.users, id)
			sm.lock.Unlock()
		}
	}
}

func (sm *SessionManager) UserConnect(uuid string) (interface{},error){
	sm.lock.Lock()
	defer sm.lock.Unlock()
	u, ok := sm.users[uuid]
	if ok {
		u.HeartBeat()
		return "ok",nil
	}
	user := new(User)
	user.lastHeartBeat = time.Now()
	user.id = uuid
	sm.users[uuid] = user
	log.Println(uuid, "connected")
	return "ok",nil
}

func (sm *SessionManager) UserDisConnect(uuid string) (interface{},error){
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if _, ok := sm.users[uuid]; ok {
		delete(sm.users, uuid)
		log.Println(uuid, "disconnected")
	}
	return "ok",nil
}

func (sm *SessionManager) GetUser(uuid string) *User {
	if u, ok := sm.users[uuid]; ok && u.IsOnline() {
		return u
	}
	return nil
}

func (sm *SessionManager) UserCount() int {
	sm.checkUser()
	return len(sm.users)
}
