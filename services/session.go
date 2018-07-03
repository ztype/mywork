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

var hbCheckTime = time.Duration(time.Second * 8)

type SessionManager struct {
	onlineusers map[string]*base.User
	lock        sync.Mutex
	db          *database.DB
	notify      chan utils.Param
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
	m.notify = make(chan utils.Param, 10)
	go m.check()
	go m.listen()
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
	return sm.notify
	return nil
}

func (sm *SessionManager) listen() {
	for {
		select {
		case p := <-sm.notify:
			if p.Type != "heartbeat" {
				sm.heartbeat(p.Uid)
			}
		}
	}
}

func (sm *SessionManager) check() {
	base.SetHeartbeatTime(hbCheckTime + time.Duration(time.Second*10))
	for {
		sm.checkUser()
		time.Sleep(hbCheckTime)
	}
}

func (sm *SessionManager) checkUser() {
	//log.Println("users", sm.onlineusers)
	for _, user := range sm.onlineusers {
		if !user.IsOnline() {
			sm.userDisConnect(user)
		}
	}
}

func (sm *SessionManager) UserConnect(uuid string) (interface{}, error) {
	u, ok := sm.onlineusers[uuid]
	var err error
	if !ok {
		u, err = sm.GetUser(uuid)
		if err == database.UserNotFound {
			u, err = sm.NewUser(uuid)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		sm.userConnect(u)
	}
	u.HeartBeat()
	return "connected", nil
}

func (sm *SessionManager) heartbeat(uuid string) {
	if u, ok := sm.onlineusers[uuid]; ok {
		u.HeartBeat()
	}
}

func (sm *SessionManager) userConnect(user *base.User) {
	if _, ok := sm.onlineusers[user.Id()]; !ok {
		sm.lock.Lock()
		defer sm.lock.Unlock()
		sm.onlineusers[user.Id()] = user
	}
}

func (sm *SessionManager) userDisConnect(user *base.User) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	uuid := user.Id()
	if _, ok := sm.onlineusers[uuid]; ok {
		delete(sm.onlineusers, uuid)
		log.Println("user", uuid, "delete")
	}
}

func (sm *SessionManager) NewUser(uuid string) (*base.User, error) {
	user := base.NewUser(uuid)
	user.HeartBeat()
	if err := sm.db.InsertUser(user); err != nil {
		return nil, err
	}
	log.Println("new user", uuid)
	return user, nil
}

func (sm *SessionManager) UserDisConnect(uuid string) (interface{}, error) {
	return "you disconnected", nil
}

func (sm *SessionManager) GetUser(uuid string) (*base.User, error) {
	return sm.db.GetUserById(uuid)
}

func (sm *SessionManager) UserCount() int {
	return len(sm.onlineusers)
}
