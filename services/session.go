package services

import (
	"log"
	"mywork/base"
	"mywork/utils"
	"sync"
	"time"
)

const (
	heartbeat = "heartbeat"
	logout    = "logout"
)

var hbCheckTime = time.Duration(time.Second * 3)

type SessionManager struct {
	users map[string]*base.User
	lock  sync.Mutex
}

type Session struct {
	user *base.User
}

func NewManager() *SessionManager {
	m := new(SessionManager)
	m.users = make(map[string]*base.User, 0)
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

func (sm *SessionManager) ObserveChannel() chan<- utils.Param {
	return nil
}

func (sm *SessionManager) check() {
	for {
		sm.checkUser()
		time.Sleep(hbCheckTime)
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

func (sm *SessionManager) UserConnect(uuid string) (interface{}, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	u, ok := sm.users[uuid]
	if ok {
		u.HeartBeat()
		return "ok", nil
	}
	user := base.NewUser(uuid)
	user.HeartBeat()
	sm.users[uuid] = user
	log.Println(uuid, "connected")
	return "ok", nil
}

func (sm *SessionManager) UserDisConnect(uuid string) (interface{}, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if _, ok := sm.users[uuid]; ok {
		delete(sm.users, uuid)
		log.Println(uuid, "disconnected")
	}
	return "ok", nil
}

func (sm *SessionManager) GetUser(uuid string) *base.User {
	if u, ok := sm.users[uuid]; ok && u.IsOnline() {
		return u
	}
	return nil
}

func (sm *SessionManager) UserCount() int {
	sm.checkUser()
	return len(sm.users)
}
