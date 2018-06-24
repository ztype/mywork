package services

import (
	"fmt"
	"log"
	"mywork/base"
	"mywork/database"
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
	onlineusers map[string]*base.User
	lock        sync.Mutex
	db          *database.DB
}

type Session struct {
	user *base.User
}

func NewManager() *SessionManager {
	m := new(SessionManager)
	m.onlineusers = make(map[string]*base.User, 0)
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	m.db = db
	go m.check()
	return m
}

func (sm *SessionManager) Name() string {
	return "connect"
}

func (sm *SessionManager) Serve(p utils.Param) (interface{}, error) {
	switch p.Type {
	case heartbeat:
		return sm.UserConnect(p.Uid)
	case logout:
		return sm.UserDisConnect(p.Uid)
	}
	return nil, fmt.Errorf("[%s] not found in %s", p.Type, sm.Name())
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
	for _, user := range sm.onlineusers {
		if !user.IsOnline() {
			sm.lock.Lock()
			sm.UserDisConnect(user.Id())
			delete(sm.onlineusers, user.Id())
			sm.lock.Unlock()
		}
	}
}

func (sm *SessionManager) UserConnect(uuid string) (interface{}, error) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	u, ok := sm.onlineusers[uuid]
	if ok {
		u.HeartBeat()
		return "ok", nil
	}
	user, err := sm.GetUser(uuid)
	if err != database.UserNotFound && err != nil {
		return nil, err
	}
	if err == database.UserNotFound {
		user, err = sm.NewUser(uuid)
		if err != nil {
			return nil, err
		}
	}
	sm.onlineusers[uuid] = user
	log.Println(uuid, "connected")
	return "ok", nil
}

func (sm *SessionManager) NewUser(uuid string) (*base.User, error) {
	user := base.NewUser(uuid)
	user.HeartBeat()
	if err := sm.db.InsertUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (sm *SessionManager) UserDisConnect(uuid string) (interface{}, error) {
	sm.db.UpdateUserOnline(uuid, false)
	log.Println(uuid, "disconnected")
	return "ok", nil
}

func (sm *SessionManager) GetUser(uuid string) (*base.User, error) {
	return sm.db.GetUserById(uuid)
}

func (sm *SessionManager) UserCount() int {
	sm.checkUser()
	return len(sm.onlineusers)
}
