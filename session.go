package main

import (
	"sync"
	"time"
)

type SessionManager struct {
	users map[string]interface{}
	lock  sync.Mutex
}

func NewManager() *SessionManager {
	m := new(SessionManager)
	m.users = make(map[string]interface{}, 0)
	return m
}

func (sm *SessionManager) UserConnect(uuid string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	_, ok := sm.users[uuid]
	if ok {
		return
	}
	sm.users[uuid] = time.Now()
}

func (sm *SessionManager) IsExist(uuid string) bool {
	_, ok := sm.users[uuid]
	return ok
}

func (sm *SessionManager) UserCount() int {
	return len(sm.users)
}
