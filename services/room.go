package services

import (
	"fmt"
	"mywork/base"
	"sync"
)

const (
	JoinRoom   = "join_room"
	CreateRoom = "create_room"
	LeaveRoom  = "leave_room"
)

type RoomId int

var idgen = []int{}

type Room struct {
	lock  sync.Mutex
	Id    RoomId
	Users map[string]*base.User
}

func (i RoomId) String() string {
	return fmt.Sprintf("%d", i)
}

func newroomid() RoomId {
	i := 1
	for i = 1; ; i++ {
		found := false
		for _, s := range idgen {
			if s == i {
				found = true
				break
			}
		}
		if !found {
			idgen = append(idgen, i)
			return RoomId(i)
		}
	}
	return RoomId(i)
}

//========= room area =============//
func newRoom() *Room {
	i := newroomid()
	r := new(Room)
	r.Id = RoomId(i)
	r.Users = make(map[string]*base.User, 0)
	return r
}

func (r *Room) ID() RoomId {
	return r.Id
}

func (r *Room) UserIn(user *base.User) error {
	ul := r.UserCount()
	if ul >= 4 {
		return fmt.Errorf("room [%s] is full", r.Id.String())
	}
	if _, ok := r.Users[user.Id()]; ok {
		return nil
	}
	r.Users[user.Id()] = user
	return nil
}

func (r *Room) UserCount() int {
	return len(r.Users)
}

func (r *Room) UserOut(user *base.User) error {
	if r.UserCount() == 0 {
		return nil
	}
	u, ok := r.Users[user.Id()]
	if ok {
		delete(r.Users, u.Id())
	}
	return nil
}

func (r *Room) HasUser(uid string) bool {
	return r.GetUser(uid) != nil
}

func (r *Room) GetUser(uid string) *base.User {
	for id, u := range r.Users {
		if id == uid {
			return u
		}
	}
	return nil
}
